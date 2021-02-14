package main

import (
	"flag"
	"github.com/sermojohn/postgres-client/pkg/config"
	"github.com/sermojohn/postgres-client/pkg/parser"
	"github.com/sermojohn/postgres-client/pkg/storage"
	"github.com/sermojohn/postgres-client/pkg/synchroniser"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	var fileName string
	flag.StringVar(&fileName, "config", "", "config file to parse")
	flag.Parse()
	if fileName == "" {
		log.Fatal("Please provide config file by using -f option")
	}

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error reading YAML file: %s\n", err)
	}

	var conf config.Config
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("Error parsing YAML file: %s\n", err)
	}

	if valid, msg := conf.Valid(); !valid {
		log.Fatalf("invalid configuration: %s", msg)
	}

	store, err := storage.New(conf.DBConfig)
	if err != nil || store == nil {
		log.Fatal(err)
	}

	synchroniser.New(store, parser.New(&http.Client{}, conf.ParserConfig)).UpdateSummaries()
}
