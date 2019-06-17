package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	db *sqlx.DB
)

func connectToDB() {
	user := os.Getenv("SQUADXML_DB_USER")
	password := os.Getenv("SQUADXML_DB_PASSWORD")
	host := os.Getenv("SQUADXML_DB_HOST")
	database := os.Getenv("SQUADXML_DB")
	var err error
	db, err = sqlx.Connect("mysql", mysqlConnection(user, password, host, database))

	if err != nil {
		logrus.WithError(err).Fatal("Error connecting to database")
	}

	if err := db.Ping(); err != nil {
		logrus.WithError(err).Fatal("Error pinging database")
	}
}

func mysqlConnection(user, password, host, database string) string {
	return user + ":" + password + "@" + host + "/" + database
}

func main() {
	connectToDB()
}
