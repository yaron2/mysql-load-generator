#!/bin/bash
set -x

. vars.sh

cat << EOF > config.json
{
    "connectionString": "${USER}:${PASSWORD}@tcp(${HOST}:${PORT})/${DB}",
    "requestsPerSecond": ${RPS},
    "printLogs": true,
    "timeToRunInSeconds": ${TTR},
    "queries" : [
	"${QUERY}"
    ]
}
EOF
