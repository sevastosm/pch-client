FROM golang:1.12-alpine

RUN apk add --no-cache git
RUN apk add --no-cache --upgrade bash

WORKDIR /parser

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["./scripts/entry.sh"]
