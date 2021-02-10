FROM golang:1.12-alpine

RUN apk add --no-cache git
RUN apk add --no-cache --upgrade bash

WORKDIR /pch-client

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

CMD ["./scripts/entry.sh"]
