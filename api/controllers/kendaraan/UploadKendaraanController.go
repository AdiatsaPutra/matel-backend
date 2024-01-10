package kendaraan

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	config "matel/configs"
	"matel/exceptions"
	"matel/models"
	"matel/payloads"
	"math"
	"mime/multipart"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbMaxIdleConns = 4
	dbMaxConns     = 100
	totalWorker    = 4
	header         = []string{
		"cabangData",
		"noKontrak",
		"namaDebitur",
		"nomorPolisi",
		"sisaHutang",
		"tipe",
		"tahun",
		"noRangka",
		"noMesin",
		"pic",
	}
)

func AddCSVPerCabang(c *gin.Context) {
	start := time.Now()

	db, err := openDbConnection(c)
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	csvReader, csvFile, err := openCsvFile(c)
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	defer func() {
		if err := csvFile.Close(); err != nil {
			exceptions.AppException(c, err.Error())
		}
	}()

	// increment cabang versi
	cabangName := c.PostForm("cabang_name")

	var cabang models.Cabang
	result := config.InitDB().Where("nama_cabang = ? AND deleted_at IS NULL", cabangName).Find(&cabang)
	if result.Error != nil {
		payloads.HandleSuccess(c, "Leasing not found", "Success", 200)
		return
	}

	cabang.Versi = cabang.Versi + 1

	result = config.InitDB().Save(&cabang)
	if result.Error != nil {
		exceptions.AppException(c, result.Error.Error())
		return
	}

	jobs := make(chan [][]interface{}, 0)
	wg := new(sync.WaitGroup)

	go dispatchWorkers(c, db, jobs, wg, cabang)
	readCsvFilePerLineThenSendToWorker(csvReader, jobs, wg)

	wg.Wait()

	duration := time.Since(start)

	var count int64
	if err := config.InitDB().Model(&models.Kendaraan{}).Count(&count).Error; err != nil {
		exceptions.AppException(c, "Something went wrong")
	}

	if err := config.InitDB().Model(&models.Home{}).Where("id = ?", 1).Update("kendaraan_total", count).Error; err != nil {
		exceptions.AppException(c, "Something went wrong")
	}

	config.CloseDB(config.InitDB())

	payloads.HandleSuccess(c, int(math.Ceil(duration.Seconds())), "Success", 200)
}

func openDbConnection(c *gin.Context) (*sql.DB, error) {
	dbConnString := ""

	if os.Getenv("GIN_MODE") == "release" {
		dbConnString = "Beta:BetaTaurus@tcp(db)/matel?charset=utf8mb4&parseTime=True&loc=Local"
	} else {
		dbConnString = "root:1Ultramilk!@tcp(127.0.0.1:3306)/motor?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := sql.Open("mysql", dbConnString)
	if err != nil {
		exceptions.AppException(c, err.Error())
		return nil, err
	}

	db.SetMaxOpenConns(dbMaxConns)
	db.SetMaxIdleConns(dbMaxIdleConns)

	return db, nil
}

func openCsvFile(c *gin.Context) (*csv.Reader, multipart.File, error) {
	file, err := c.FormFile("file")
	if err != nil {
		exceptions.AppException(c, err.Error())
		return nil, nil, err
	}

	csvFile, err := file.Open()
	if err != nil {
		exceptions.AppException(c, err.Error())
		return nil, nil, err
	}

	sniffer := make([]byte, 4096)
	_, err = csvFile.Read(sniffer)
	if err != nil {
		exceptions.AppException(c, err.Error())
		return nil, nil, err
	}

	_, err = csvFile.Seek(0, io.SeekStart)
	if err != nil {
		exceptions.AppException(c, err.Error())
		return nil, nil, err
	}

	delimiter := detectDelimiter(sniffer)

	reader := csv.NewReader(csvFile)
	reader.Comma = delimiter

	return reader, csvFile, nil
}

func detectDelimiter(content []byte) rune {
	if bytes.Contains(content, []byte(";")) {
		return ';'
	} else if bytes.Contains(content, []byte("\t")) {
		return '\t'
	}

	return ','
}

func dispatchWorkers(c *gin.Context, db *sql.DB, jobs <-chan [][]interface{}, wg *sync.WaitGroup, cabang models.Cabang) {
	for workerIndex := 0; workerIndex <= totalWorker; workerIndex++ {
		go func(workerIndex int, db *sql.DB, jobs <-chan [][]interface{}, wg *sync.WaitGroup) {
			for job := range jobs {
				err := doTheJobBatch(c, workerIndex, db, job, cabang)
				if err != nil {
					exceptions.AppException(c, err.Error())
				}
				wg.Done()
			}
		}(workerIndex, db, jobs, wg)
	}
}

func readCsvFilePerLineThenSendToWorker(csvReader *csv.Reader, jobs chan<- [][]interface{}, wg *sync.WaitGroup) {
	isHeader := true
	batchSize := 100
	batch := make([][]interface{}, 0, batchSize)

	for {
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}

		if isHeader {
			isHeader = false
			continue
		}

		rowOrdered := make([]interface{}, 0)
		for _, each := range row {
			rowOrdered = append(rowOrdered, each)
		}

		batch = append(batch, rowOrdered)

		if len(batch) >= batchSize {
			wg.Add(1)
			jobs <- batch
			batch = make([][]interface{}, 0, batchSize)
		}
	}

	if len(batch) > 0 {
		wg.Add(1)
		jobs <- batch
	}

	close(jobs)
}

func doTheJobBatch(c *gin.Context, workerIndex int, db *sql.DB, rows [][]interface{}, cabang models.Cabang) error {
	now := time.Now()

	leasingName := c.PostForm("leasing_name")
	cabangIDStr := c.PostForm("cabang_id")
	cabangID, err := strconv.Atoi(cabangIDStr)
	if err != nil {
		exceptions.AppException(c, "Invalid cabang_id")
		return err
	}

	cabangName := c.PostForm("cabang_name")

	var valuesBatch []interface{}
	for _, row := range rows {
		values := make([]interface{}, 0)
		values = append(values, cabangID)
		values = append(values, cabangName)
		values = append(values, row...)

		values = append([]interface{}{leasingName}, values...)
		values = append(values, now)
		values = append(values, 1)
		values = append(values, cabang.Versi)

		var alphanumericRegex = regexp.MustCompile("[^a-zA-Z0-9]+")

		for i := 6; i < 9; i++ {
			if str, ok := values[i].(string); ok {
				filteredStr := alphanumericRegex.ReplaceAllString(str, "")
				values[i] = filteredStr
			}
		}

		for i := 10; i < 12; i++ {
			if str, ok := values[i].(string); ok {
				filteredStr := alphanumericRegex.ReplaceAllString(str, "")
				values[i] = filteredStr
			}
		}

		valuesBatch = append(valuesBatch, values)
	}

	conn, err := db.Conn(context.Background())
	if err != nil {
		exceptions.AppException(c, err.Error())
		return err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			exceptions.AppException(c, err.Error())
		}
	}()

	placeholderRows := generateQuestionsMark(len(rows), len(header)+6)

	query := fmt.Sprintf("INSERT INTO m_kendaraan (leasing, cabang_id, cabang, %s, created_at, status, versi) VALUES %s",
		strings.Join(header, ","),
		placeholderRows,
	)

	args := make([]interface{}, 0, len(valuesBatch)*len(header)+6)
	for _, values := range valuesBatch {
		if v, ok := values.([]interface{}); ok {
			args = append(args, v...)
		}
	}

	_, err = conn.ExecContext(context.Background(), query, args...)
	if err != nil {
		exceptions.AppException(c, err.Error())
		return err
	}

	return nil
}

func generateQuestionsMark(rowsCount, columnsCount int) string {
	placeholder := "(?" + strings.Repeat(", ?", columnsCount-1) + ")"
	return strings.Repeat(", "+placeholder, rowsCount)[2:]
}
