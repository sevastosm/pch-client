package main

import (
	"flag"
	"github.com/sermojohn/postgres-client/pkg/config"
	"github.com/sermojohn/postgres-client/pkg/parser"
	"github.com/sermojohn/postgres-client/pkg/storage"
	"github.com/sermojohn/postgres-client/pkg/synchroniser"
	"log"
	"net/http"
	"os"
)

func main() {
	dbConfig := config.DBConfig{}
	clientConfig := config.ClientConfig{}
	flag.StringVar(&dbConfig.Host, "db.host", "localhost", "database host")
	flag.IntVar(&dbConfig.Port, "db.port", 5432, "database port")
	flag.StringVar(&dbConfig.User, "db.user", "ioannis", "database user to connect with")
	flag.StringVar(&dbConfig.Password, "db.password", "ioannissecret", "database password of specified user")
	flag.StringVar(&dbConfig.Name, "db.name", "ioannisdb", "database name")
	flag.Int64Var(&clientConfig.ParserRateLimitDelayMillis, "rate-limit-delay", 0, "parser delay between query request")
	flag.IntVar(&clientConfig.ServerLimit, "amount-of-servers", 10, "the amount of servers to query")
	flag.StringVar(&clientConfig.IPVersion, "ip-version", "v4", "IP version")
	flag.StringVar(&clientConfig.IXP, "ixp", "", "IXP ID")
	flag.StringVar(&clientConfig.City, "city", "", "IXP server city")
	flag.StringVar(&clientConfig.Country, "country", "", "IXP server country")
	flag.Parse()

	if dbConfig.HasMissingConfig() {
		flag.Usage()
		os.Exit(1)
	}

	store, err := storage.New(dbConfig)
	if err != nil || store == nil {
		log.Fatal(err)
	}

	synchroniser.New(store, parser.New(&http.Client{}, clientConfig)).UpdateSummaries()
}
