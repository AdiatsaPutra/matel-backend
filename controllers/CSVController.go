package controllers

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"motor/exceptions"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var (
	// dbConnString = "root:1Ultramilk!@tcp(127.0.0.1:3306)/motor?charset=utf8mb4&parseTime=True&loc=Local"
	dbConnString   = "w08um7qaben07grspf9k:pscale_pw_VAkTxIR732WX6GQhmtlAamddhm7CSHSHhY69U2rjIm7@tcp(aws.connect.psdb.cloud)/matel?tls=true&charset=utf8mb4&parseTime=True&loc=Local"
	dbMaxIdleConns = 4
	dbMaxConns     = 100
	totalWorker    = 100
	// csvFile        = "majestic_million.csv"
	dataHeaders = []string{
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
		"status",
	}
)

func openDbConnection() (*sql.DB, error) {
	log.Println("=> open db connection")

	db, err := sql.Open("mysql", dbConnString)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(dbMaxConns)
	db.SetMaxIdleConns(dbMaxIdleConns)

	return db, nil
}

func AddCSV(c *gin.Context) {
	start := time.Now()
	csvReader, csvFile, err := openCsvFile(c)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer csvFile.Close()

	jobs := make(chan []interface{}, 0)
	wg := new(sync.WaitGroup)

	db, err := openDbConnection()
	if err != nil {
		log.Fatal(err.Error())
	}

	go dispatchWorkers(db, jobs, wg)
	readCsvFilePerLineThenSendToWorker(csvReader, jobs, wg)

	wg.Wait()

	duration := time.Since(start)
	log.Println(fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds"))

}

func openCsvFile(c *gin.Context) (*csv.Reader, multipart.File, error) {

	// Retrieve the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		logrus.Info(err)
		exceptions.BadRequest(c, "Masukkan data valid")
		return nil, nil, err
	}
	logrus.Info(file)

	// Open the uploaded file
	csvFile, err := file.Open()
	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return nil, nil, err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	return reader, csvFile, nil
}

func dispatchWorkers(db *sql.DB, jobs <-chan []interface{}, wg *sync.WaitGroup) {
	for workerIndex := 0; workerIndex <= totalWorker; workerIndex++ {
		go func(workerIndex int, db *sql.DB, jobs <-chan []interface{}, wg *sync.WaitGroup) {
			counter := 0

			for job := range jobs {
				doTheJob(workerIndex, counter, db, job)
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

func doTheJob(workerIndex, counter int, db *sql.DB, values []interface{}) {
	for {
		var outerError error
		func(outerError *error) {
			defer func() {
				if err := recover(); err != nil {
					*outerError = fmt.Errorf("%v", err)
				}
			}()

			conn, err := db.Conn(context.Background())
			// logrus.Info(err)
			query := fmt.Sprintf("INSERT INTO m_leasing (%s) VALUES (%s)",
				strings.Join(dataHeaders, ","),
				strings.Join(generateQuestionsMark(len(dataHeaders)), ","),
			)

			logrus.Info(query)
			_, err = conn.ExecContext(context.Background(), query, values...)
			// logrus.Info("INSERT")
			if err != nil {
				logrus.Info(query)
				log.Println(err)
				log.Fatal(err.Error())
			}

			err = conn.Close()
			if err != nil {
				log.Fatal(err.Error())
			}
		}(&outerError)
		if outerError == nil {
			break
		}
	}

	// if counter%100 == 0 {
	// 	log.Println("=> worker", workerIndex, "inserted", counter, "data")
	// }
}

func generateQuestionsMark(n int) []string {
	s := make([]string, 0)
	for i := 0; i < n; i++ {
		s = append(s, "?")
	}
	return s
}
