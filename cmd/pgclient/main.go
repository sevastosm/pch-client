package main

import (
	"flag"
	"fmt"
	"github.com/sermojohn/postgres-client/pkg/config"
	"github.com/sermojohn/postgres-client/pkg/user"
	"os"
)

func main() {
	dbConfig := config.DBConfig{}
	flag.StringVar(&dbConfig.Host, "host", "", "database host")
	flag.IntVar(&dbConfig.Port, "port", 0, "database port")
	flag.StringVar(&dbConfig.User, "user", "", "database user to connect with")
	flag.StringVar(&dbConfig.Password, "password", "", "database password of specified user")
	flag.StringVar(&dbConfig.Name, "name", "", "database name")
	flag.Parse()

	if dbConfig.HasMissingConfig() {
		flag.Usage()
		os.Exit(1)
	}

	userStore, err := user.New(dbConfig)
	if err != nil {
		fmt.Printf("client failed with: %v", err)
		os.Exit(1)
	}

	user, found := userStore.FindUser(5)
	if !found {
		fmt.Printf("queried user not found")
		os.Exit(1)
	}
	fmt.Printf("query returned user: %v\n", user)

	allUsers, err := userStore.FindAllUsers()
	if err != nil {
		fmt.Printf("query failed with: %v", err)
		os.Exit(1)
	}
	for _, u := range allUsers {
		fmt.Printf("user row: %v\n", u)
	}
}


