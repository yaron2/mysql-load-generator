package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	RequestsPerSecond  int      `json:"requestsPerSecond"`
	Queries            []string `json:"queries"`
	ConnectionString   string   `json:"connectionString"`
	PrintLogs          bool     `json:"printLogs"`
	TimeToRunInSeconds int      `json:"timeToRunInSeconds"`
}

var config Config
var db *sql.DB
var queryTimesInMS []int

var wg sync.WaitGroup

func main() {
	loadConfig()
	openMemSQLConnection()
	dispatchQueries()
	displayAverageQueryTime()
}

func openMemSQLConnection() {
	var err error
	db, err = sql.Open("mysql", config.ConnectionString)

	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connection succeeded")
}

func trackFuncTime(start time.Time) {
	elapsed := time.Since(start)
	queryTimesInMS = append(queryTimesInMS, int(elapsed/time.Millisecond))

	if config.PrintLogs {
		fmt.Println(elapsed)
	}
}

func executeQuery(query string) {
	defer wg.Done()
	defer trackFuncTime(time.Now())

	_, err := db.Exec(query)

	if err != nil {
		fmt.Println(err.Error())
	}

	if config.PrintLogs {
		fmt.Println("Executed " + query)
	}
}

func displayAverageQueryTime() {
	if config.PrintLogs {
		var totalvalues int

		for _, val := range queryTimesInMS {
			totalvalues += val
		}

		var count = len(queryTimesInMS)
		if count != 0{
			sum := totalvalues / count
			fmt.Println("Average time: " + strconv.Itoa(sum) + "ms")
		}
	}
}

func dispatchQueries() {
	fmt.Println("Starting run")

	start := time.Now()
	shouldEnd := start.Add(time.Second * time.Duration(config.TimeToRunInSeconds))

	for time.Now().Before(shouldEnd) {
		for i := 0; i < config.RequestsPerSecond; i++ {
			wg.Add(1)
			go executeQuery(config.Queries[rand.Intn(len(config.Queries))])
		}

		time.Sleep(1 * time.Second)
	}

	wg.Wait()
}

func loadConfig() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	json.Unmarshal(data, &config)
	fmt.Println("Config loaded")
}
