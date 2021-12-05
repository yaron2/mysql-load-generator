package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/namsral/flag"

	_ "github.com/go-sql-driver/mysql"
)

type Substitution struct {
	Key string `json:"key"`
	Min int    `json:"min"`
	Max int    `json:"max"`
}

type Config struct {
	RequestsPerSecond  int            `json:"requestsPerSecond"`
	Queries            []string       `json:"queries"`
	ConnectionString   string         `json:"connectionString"`
	PrintLogs          bool           `json:"printLogs"`
	TimeToRunInSeconds int            `json:"timeToRunInSeconds"`
	PoolConnections    int            `json:"poolConnections"`
	DryRun             bool           `json:"dryRun"`
	QueryTimeout       int            `json:"queryTimeout"`
	ConnectionLifetime time.Duration  `json:"connectionLifeTime"`
	Substitution       []Substitution `json:"substitution"`
}

var config = Config{
	RequestsPerSecond:  100,
	Queries:            make([]string, 0),
	ConnectionString:   "",
	PrintLogs:          true,
	TimeToRunInSeconds: 10,
	PoolConnections:    10,
	DryRun:             true,
	QueryTimeout:       10,
	ConnectionLifetime: 0,
	Substitution:       []Substitution{},
}
var configFilePath string
var db *sql.DB
var queryTimesInMS []int
var ops int64 = 0
var latency int64 = 0

var wg sync.WaitGroup

func main() {
	flag.StringVar(&configFilePath, "configFile", "/app/config.json", "path to config file")
	flag.Parse()
	go printQps()
	loadConfig(configFilePath)
	openMemSQLConnection()
	dispatchQueries()
}

func printQps() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		loadedOps := atomic.LoadInt64(&ops)
		atomic.StoreInt64(&ops, 0)
		loadedLatency := atomic.LoadInt64(&latency)
		atomic.StoreInt64(&latency, 0)
        if loadedOps != 0 {
		    fmt.Println("qps: ", loadedOps/5, ", average latency: ", time.Duration(loadedLatency/loadedOps))
        } else {
            fmt.Println("qps: ", loadedOps/5)
        }

	}

}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

func openMemSQLConnection() {
	var err error
	db, err = sql.Open("mysql", config.ConnectionString)
	db.SetConnMaxLifetime(config.ConnectionLifetime)
	db.SetMaxIdleConns(config.PoolConnections)
	db.SetMaxOpenConns(config.PoolConnections)

	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connection succeeded with pool size", config.PoolConnections)
}

func executeQuery(query string) {
	defer wg.Done()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.QueryTimeout)*time.Second)
	defer cancel()

	for _, subs := range config.Substitution {
		query = strings.Replace(query, subs.Key, strconv.Itoa(rand.Intn(subs.Max-subs.Min)+subs.Min), 1)
	}

	if config.PrintLogs {
		fmt.Println("Executed ", query)
	}
	if config.DryRun == false {
		start := time.Now()
		_, err := db.ExecContext(ctx, query)
		if err == nil {
			atomic.AddInt64(&ops, 1)
			atomic.AddInt64(&latency, int64(time.Since(start)))
		} else {
			fmt.Fprintf(os.Stderr, "query err: %s\n", err.Error())
		}
	}
}

func dispatchQueries() {
	fmt.Println("Starting run")

	start := time.Now()
	shouldEnd := start.Add(time.Second * time.Duration(config.TimeToRunInSeconds))

	for time.Now().Before(shouldEnd) {
		l := len(config.Queries)
		for i := 0; i < config.RequestsPerSecond/l; i++ {
			wg.Add(l)
			for q := 0; q < l; q++ {
				go executeQuery(config.Queries[q])
			}
		}

		time.Sleep(1 * time.Second)
	}
	fmt.Println("wait")
	if waitTimeout(&wg, 10*time.Second) {
		fmt.Println("Timed out waiting for wait group")
	} else {
		fmt.Println("Wait group finished")
	}
	fmt.Println("end wait")
	defer db.Close()
}

func loadConfig(configFilePath string) {
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		panic(err.Error() + " use -configFile flag")
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	fmt.Println("LOADED CONFIG")
	pretty, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(pretty))

	if len(config.Queries) < 1 {
		panic("no query specified")
	}
}
