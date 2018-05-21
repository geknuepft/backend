package database

import (
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
	"os"
)

var dbPoolHandle *sqlx.DB

func init() {
	var err error

	mysql_host := os.Getenv("MYSQL_HOST")

	dbPoolHandle, err = sqlx.Connect(
		"mysql",
		"geknuepft:Er3cof4iesho@tcp("+mysql_host+":3306)/geknuepft",
	)
	if err != nil {
		panic(err.Error())
	}
}

func GetDbX() (db *sqlx.DB) {
	return dbPoolHandle
}
