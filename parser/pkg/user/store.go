package user

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sermojohn/postgres-client/pkg/config"
)

// Store abstraction for user data storage
type Store interface {
	CreateUser(id int, username string) error
	FindUser(id int) (*User, bool)
	FindAllUsers() ([]User, error)
}

type userStore struct {
	db *sql.DB
}

func (us *userStore) CreateUser(id int, username string) error {
	_, err := us.db.Exec("insert into Users values ($1, $2);", id, username)
	if err != nil {
		return err
	}
	return nil
}

func (us *userStore) FindUser(id int) (*User, bool) {
	var u User

	err := us.db.QueryRow("select * from Users where id = $1", id).Scan(&u.id, &u.username)
	if err != nil {
		fmt.Printf("%v", err)
		return nil, false
	}

	return &u, true
}

func (us *userStore) FindAllUsers() ([]User, error) {
	var users []User

	rows, err := us.db.Query("select * from Users order by id")
	if err != nil {
		return nil, err
	}

	for ; rows.Next(); {
		nu := User{}
		err := rows.Scan(&nu.id, &nu.username)
		if err != nil {
			return nil, err
		}
		users = append(users, nu)
	}

	return users, nil
}

func New(dbc config.DBConfig) (Store, error) {
	pgConn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbc.Host, dbc.Port, dbc.User, dbc.Password, dbc.Name)

	db, err := sql.Open("postgres", pgConn)
	if err != nil {
		return nil, err
	}

	return &userStore{db: db}, nil
}
