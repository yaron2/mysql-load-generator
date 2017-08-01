# MemSQL / MySQL Load Tester

This is a tool for generating load tests based on custom queries for MemSQL / MySQL databases.

## How does it work?

The configuration resides in config.json.
Queries will be randomly executed against the target Database.

Json configuration guide:

* connectionString (string) - the connection string to the MySQL / MemSQL Database
* requestsPerSecond (int) - how many statements per second to execute
* printLogs (bool) - print queries executed including time taken for each query and a total average
* timeToRun (int) - how long to run the test in seconds
* queries (string array) - the queries to run

## Run as Docker

Simply build using the Dockerfile.
Load tests are a great fit for [Azure Container Instances].

[Azure Container Instances]: https://azure.microsoft.com/en-us/services/container-instances/