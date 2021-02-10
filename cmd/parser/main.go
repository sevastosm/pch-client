package main

import (
	"flag"
	"github.com/sermojohn/postgres-client/pkg/config"
	"github.com/sermojohn/postgres-client/pkg/parser"
	"github.com/sermojohn/postgres-client/pkg/storage"
	"github.com/sermojohn/postgres-client/pkg/updater"
	"log"
	"net/http"
	"os"
)

func main() {
	dbConfig := config.DBConfig{}
	updaterConfig := config.UpdateConfig{}
	parserConfig := config.ParserConfig{}
	flag.StringVar(&dbConfig.Host, "db.host", "localhost", "database host")
	flag.IntVar(&dbConfig.Port, "db.port", 5432, "database port")
	flag.StringVar(&dbConfig.User, "db.user", "ioannis", "database user to connect with")
	flag.StringVar(&dbConfig.Password, "db.password", "ioannissecret", "database password of specified user")
	flag.StringVar(&dbConfig.Name, "db.name", "ioannisdb", "database name")
	flag.Int64Var(&updaterConfig.ParserRateLimitDelayMillis, "rate-limit-delay", 0, "parser delay between query request")
	flag.IntVar(&updaterConfig.AmountOfServers, "amount-of-servers", 10, "the amount of servers to query")
	flag.StringVar(&parserConfig.IPVersion, "ip-version", "v4", "IP version")
	flag.Parse()

	if dbConfig.HasMissingConfig() {
		flag.Usage()
		os.Exit(1)
	}

	store, err := storage.New(dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	sync := updater.New(store, parser.New(&http.Client{}, parserConfig), updaterConfig)
	err = sync.UpdateSummaries()
	if err != nil {
		log.Printf("failed to update summaries with error %v", err)
	}
}
