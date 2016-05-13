package main

import (
	"database/sql"

	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

// Give us some seed data
func init() {
	var err error

	Db, err = sql.Open("mysql", "geknuepft:Er3cof4iesho@tcp(dc_mysql-server.docker.:3306)/geknuepft")
	if err != nil {
		panic(err.Error())
	}
	//defer Db.Close()

	err = Db.Ping()
	if err != nil {
		panic(err.Error())
	}
}

type rawTime []byte

func (t rawTime) Time() time.Time {
	time, err := time.Parse("2006-01-02 15:04:05", string(t))
	if err != nil {
		panic(err.Error())
	}
	return time
}

func GetArticles() (articles Articles) {
	rows, err := Db.Query(
		`SELECT article_id,
		COALESCE(article_name_de, ''),
		created
		FROM article
		WHERE article_id`)

	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var a Article

		var created rawTime

		err := rows.Scan(&a.Id, &a.Name, &created)
		if err != nil {
			panic(err.Error())
		}

		a.Created = created.Time()

		articles = append(articles, a)
	}

	return
}
