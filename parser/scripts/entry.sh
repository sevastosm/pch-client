#!/bin/sh

trap "echo TRAPed signal" HUP INT QUIT TERM

CMD=${PARSER:-"go run ./cmd/parser/main.go"}
PARAMS="-db.host=${DB_HOST:-localhost}"
INTERVAL=$1
INTERVAL=${INTERVAL:-10}

# run job
while true; do ($CMD $PARAMS); sleep $INTERVAL; done
