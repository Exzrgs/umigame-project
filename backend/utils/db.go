package utils

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     string
	dbDatabase string
	dbConfig   string
)

func ConnectDB() (*sqlx.DB, error) {
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	dbDatabase = os.Getenv("DB_DATABASE")
	dbConfig = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbDatabase)

	db, err := sqlx.Open("mysql", dbConfig)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
