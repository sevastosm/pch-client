#!/bin/sh

trap exit HUP INT QUIT TERM

CMD="go run ./cmd/parser/main.go"
PARAMS="-config=/config/parser.yml"
INTERVAL=${PARSER_RUN_EVERY_SEC:-10}

# run job
while true; do ($CMD $PARAMS); sleep $INTERVAL; done
