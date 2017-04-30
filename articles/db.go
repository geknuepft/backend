package articles

import (
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

func getDbX() (db *sqlx.DB) {
	var err error

	db, err = sqlx.Connect("mysql", "geknuepft:Er3cof4iesho@tcp(mysql-server:3306)/geknuepft")
	if err != nil {
		panic(err.Error())
	}

	return
}
