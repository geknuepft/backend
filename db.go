package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// Give us some seed data
func init() {
	db, err := sql.Open("mysql", "geknuepft:Er3cof4iesho@tcp(dc_mysql-server.docker.:3306)/geknuepft")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
}
