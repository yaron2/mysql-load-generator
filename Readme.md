# MemSQL / MySQL Load Tester

This is a tool for generating load tests based on custom queries for MemSQL / MySQL databases.
This fork brings substitution, query timeouts, wait group timeout, connection pool, periodic qps and latency logs, higher performance

## How does it work?

The configuration resides in config.json.
Queries will be randomly executed against the target Database.

Json configuration guide:

* connectionString (string) - the connection string to the MySQL / MemSQL Database
* requestsPerSecond (int) - how many statements per second to execute
* printLogs (bool) - print queries executed 
* timeToRun (int) - how long to run the test in seconds
* queries (string array) - the queries to run, can use substitution placeholder
* substitution ( map ) - specify placeholder for substitution with random number according specified min,max (int) rules

## Run as Docker
Prepare config.json and run
```
docker run --network=host -v $(pwd)/config.json:/app/config.json -it --rm tombokombo/mysql-loader:latest
```
