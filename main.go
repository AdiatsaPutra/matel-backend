package main

// import (
// 	"context"
// 	"database/sql"
// 	"encoding/csv"
// 	"fmt"
// 	"io"
// 	"log"
// 	"math"
// 	"os"
// 	"strings"
// 	"sync"
// 	"time"

// 	_ "github.com/go-sql-driver/mysql"
// )

// var (
// 	dbConnString   = "root@/test"
// 	dbMaxIdleConns = 4
// 	dbMaxConns     = 100
// 	totalWorker    = 100
// 	csvFile        = "majestic_million.csv"
// 	dataHeaders    = []string{
// 		"GlobalRank",
// 		"TldRank",
// 		"Domain",
// 		"TLD",
// 		"RefSubNets",
// 		"RefIPs",
// 		"IDN_Domain",
// 		"IDN_TLD",
// 		"PrevGlobalRank",
// 		"PrevTldRank",
// 		"PrevRefSubNets",
// 		"PrevRefIPs",
// 	}
// )

// // CREATE DATABASE IF NOT EXISTS test;
// // USE test;
// // CREATE TABLE IF NOT EXISTS domain (
// //     GlobalRank int,
// //     TldRank int,
// //     Domain varchar(255),
// //     TLD varchar(255),
// //     RefSubNets int,
// //     RefIPs int,
// //     IDN_Domain varchar(255),
// //     IDN_TLD varchar(255),
// //     PrevGlobalRank int,
// //     PrevTldRank int,
// //     PrevRefSubNets int,
// //     PrevRefIPs int
// // );

// func main() {
// 	start := time.Now()

// 	db, err := openDbConnection()
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	csvReader, csvFile, err := openCsvFile()
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// 	defer csvFile.Close()

// 	jobs := make(chan []interface{}, 0)
// 	wg := new(sync.WaitGroup)

// 	go dispatchWorkers(db, jobs, wg)
// 	readCsvFilePerLineThenSendToWorker(csvReader, jobs, wg)

// 	wg.Wait()

// 	duration := time.Since(start)
// 	fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")
// }

// func openDbConnection() (*sql.DB, error) {
// 	log.Println("=> open db connection")

// 	db, err := sql.Open("mysql", dbConnString)
// 	if err != nil {
// 		return nil, err
// 	}

// 	db.SetMaxOpenConns(dbMaxConns)
// 	db.SetMaxIdleConns(dbMaxIdleConns)

// 	return db, nil
// }

// func openCsvFile() (*csv.Reader, *os.File, error) {
// 	log.Println("=> open csv file")

// 	f, err := os.Open(csvFile)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			log.Fatal("file majestic_million.csv tidak ditemukan. silakan download terlebih dahulu di https://blog.majestic.com/development/majestic-million-csv-daily")
// 		}

// 		return nil, nil, err
// 	}

// 	reader := csv.NewReader(f)
// 	return reader, f, nil
// }

// func dispatchWorkers(db *sql.DB, jobs <-chan []interface{}, wg *sync.WaitGroup) {
// 	for workerIndex := 0; workerIndex <= totalWorker; workerIndex++ {
// 		go func(workerIndex int, db *sql.DB, jobs <-chan []interface{}, wg *sync.WaitGroup) {
// 			counter := 0

// 			for job := range jobs {
// 				doTheJob(workerIndex, counter, db, job)
// 				wg.Done()
// 				counter++
// 			}
// 		}(workerIndex, db, jobs, wg)
// 	}
// }

// func readCsvFilePerLineThenSendToWorker(csvReader *csv.Reader, jobs chan<- []interface{}, wg *sync.WaitGroup) {
// 	isHeader := true
// 	for {
// 		row, err := csvReader.Read()
// 		if err != nil {
// 			if err == io.EOF {
// 				err = nil
// 			}
// 			break
// 		}

// 		if isHeader {
// 			isHeader = false
// 			continue
// 		}

// 		rowOrdered := make([]interface{}, 0)
// 		for _, each := range row {
// 			rowOrdered = append(rowOrdered, each)
// 		}

// 		wg.Add(1)
// 		jobs <- rowOrdered
// 	}
// 	close(jobs)
// }

// func doTheJob(workerIndex, counter int, db *sql.DB, values []interface{}) {
// 	for {
// 		var outerError error
// 		func(outerError *error) {
// 			defer func() {
// 				if err := recover(); err != nil {
// 					*outerError = fmt.Errorf("%v", err)
// 				}
// 			}()

// 			conn, err := db.Conn(context.Background())
// 			query := fmt.Sprintf("INSERT INTO domain (%s) VALUES (%s)",
// 				strings.Join(dataHeaders, ","),
// 				strings.Join(generateQuestionsMark(len(dataHeaders)), ","),
// 			)

// 			_, err = conn.ExecContext(context.Background(), query, values...)
// 			if err != nil {
// 				log.Fatal(err.Error())
// 			}

// 			err = conn.Close()
// 			if err != nil {
// 				log.Fatal(err.Error())
// 			}
// 		}(&outerError)
// 		if outerError == nil {
// 			break
// 		}
// 	}

// 	if counter%100 == 0 {
// 		log.Println("=> worker", workerIndex, "inserted", counter, "data")
// 	}
// }

// func generateQuestionsMark(n int) []string {
// 	s := make([]string, 0)
// 	for i := 0; i < n; i++ {
// 		s = append(s, "?")
// 	}
// 	return s
// }

import (
	"database/sql"
	"fmt"
	"motor/controllers"
	"motor/security"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	r := gin.Default()

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	r.GET("/profil", security.AuthMiddleware(), controllers.GetProfile)

	r.GET("/leasing", controllers.GetLeasing)
	r.POST("/upload-leasing", controllers.UploadLeasing)
	r.GET("/export", exportHandler)

	r.GET("/member", security.AuthMiddleware(), controllers.GetAllMember)

	r.GET("/province", controllers.GetProvince)
	r.GET("/kabupaten/:province-id", controllers.GetKabupaten)
	r.GET("/kecamatan/:kabupaten-id", controllers.GetKecamatan)

	r.Run()
}

type Leasing struct {
	ID   int
	Name string
	// Add other fields as needed
}

func exportHandler(c *gin.Context) {
	// Retrieve all data from the table
	var data []Leasing

	// Populate data with sample data
	data = append(data, Leasing{ID: 1, Name: "Lease 1"})
	data = append(data, Leasing{ID: 2, Name: "Lease 2"})
	// Add more sample data as needed

	// Create a new SQLite database file
	sqliteDB, err := sql.Open("sqlite3", "exported.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer sqliteDB.Close()

	// Enable foreign key support in SQLite
	_, err = sqliteDB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create the Leasing table in the SQLite database
	_, err = sqliteDB.Exec(`CREATE TABLE IF NOT EXISTS leasing (
		id INTEGER PRIMARY KEY,
		name TEXT
	);`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Prepare the INSERT statement
	stmt, err := sqliteDB.Prepare(`INSERT INTO leasing (
		id, name
	) VALUES (?, ?);`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	// Insert data into the SQLite database
	for _, item := range data {
		_, err = stmt.Exec(item.ID, item.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Set the response headers for file download
	filename := "exported.db"
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.File(filepath.Join(".", filename))
}
