package controllers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"matel/exceptions"
	"matel/payloads"
	"math"
	"mime/multipart"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var (
	// dbConnString = "root:1Ultramilk!@tcp(127.0.0.1:3306)/motor?charset=utf8mb4&parseTime=True&loc=Local"
	// dbConnString = "root:root@tcp(167.172.69.241:3306)/matel?charset=utf8mb4&parseTime=True&loc=Local"
	// dbConnString = "root:root@tcp(db)/matel?charset=utf8mb4&parseTime=True&loc=Local"

	dbMaxIdleConns = 4
	dbMaxConns     = 100
	totalWorker    = 100
	dataHeaders    = []string{
		"leasing",
		"cabang",
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

func openDbConnection(c *gin.Context) (*sql.DB, error) {

	dbConnString := ""

	if os.Getenv("GIN_MODE") == "release" {
		dbConnString = "root:root@tcp(db)/matel?charset=utf8mb4&parseTime=True&loc=Local"
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

func AddCSV(c *gin.Context) {
	start := time.Now()

	db, err := openDbConnection(c)
	if err != nil {

		exceptions.AppException(c, err.Error())
		return
	}

	csvReader, _, err := openCsvFile(c)
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}
	// defer csvFile.Close()

	jobs := make(chan []interface{}, 0)
	wg := new(sync.WaitGroup)

	go dispatchWorkers(c, db, jobs, wg)
	readCsvFilePerLineThenSendToWorker(csvReader, jobs, wg)

	wg.Wait()

	duration := time.Since(start)
	payloads.HandleSuccess(c, int(math.Ceil(duration.Seconds())), "Success", 200)

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

func dispatchWorkers(c *gin.Context, db *sql.DB, jobs <-chan []interface{}, wg *sync.WaitGroup) {
	for workerIndex := 0; workerIndex <= totalWorker; workerIndex++ {
		go func(workerIndex int, db *sql.DB, jobs <-chan []interface{}, wg *sync.WaitGroup) {
			counter := 0

			for job := range jobs {
				doTheJob(c, workerIndex, counter, db, job)
				wg.Done()
				counter++
			}
		}(workerIndex, db, jobs, wg)
	}
}

func readCsvFilePerLineThenSendToWorker(csvReader *csv.Reader, jobs chan<- []interface{}, wg *sync.WaitGroup) {
	isHeader := true
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

		wg.Add(1)
		jobs <- rowOrdered
	}
	close(jobs)
}

func doTheJob(c *gin.Context, workerIndex, counter int, db *sql.DB, values []interface{}) {
	now := time.Now()

	var alphanumericRegex = regexp.MustCompile("[^a-zA-Z0-9]+")

	for i := 4; i < 7; i++ {
		if str, ok := values[i].(string); ok {
			filteredStr := alphanumericRegex.ReplaceAllString(str, "")
			values[i] = filteredStr
		}
	}

	leasingName := c.PostForm("leasing_name")
	cabangName := c.PostForm("cabang_name")

	if leasingName != "" {
		values[0] = leasingName
	}

	if cabangName != "" {
		values[1] = cabangName
	}

	values = append(values, now)
	values = append(values, 1)

	for {
		var outerError error
		func(outerError *error) {
			defer func() {
				if err := recover(); err != nil {
					*outerError = fmt.Errorf("%v Error", err)
				}
			}()

			conn, err := db.Conn(context.Background())
			query := fmt.Sprintf("INSERT INTO m_kendaraan (%s, created_at, status) VALUES (%s)",
				strings.Join(dataHeaders, ","),
				strings.Join(generateQuestionsMark(len(dataHeaders)+2), ","),
			)

			_, err = conn.ExecContext(context.Background(), query, values...)
			if err != nil {
				exceptions.AppException(c, err.Error())
				return
			}

			err = conn.Close()
			if err != nil {
					exceptions.AppException(c, err.Error())
					return
			}
		}(&outerError)
		if outerError == nil {
			break
		}
	}
}

func generateQuestionsMark(n int) []string {
	s := make([]string, 0)
	for i := 0; i < n; i++ {
		s = append(s, "?")
	}
	return s
}