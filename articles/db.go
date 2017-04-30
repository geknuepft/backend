package articles

import (
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
	"os"
)

func getDbX() (db *sqlx.DB) {
	var err error

	mysql_host := os.Getenv("MYSQL_HOST")

	db, err = sqlx.Connect(
		"mysql",
		"geknuepft:Er3cof4iesho@tcp("+mysql_host+":3306)/geknuepft",
 	)
	if err != nil {
		panic(err.Error())
	}

	return
}
