package main

import (
	"database/sql"
	"strings"

	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

// Give us some seed data
func init() {
	var err error

	Db, err = sql.Open("mysql", "geknuepft:Er3cof4iesho@tcp(172.17.0.2:3306)/geknuepft")
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
		`SELECT a.article_id,
		COALESCE(a.article_name_de, ''),
		a.created,
		i0.path p0,
		i1.path p1
		FROM article a
		JOIN image_type it0 ON(it0.abbr = 'rma0')
		LEFT JOIN image i0 ON(i0.article_id = a.article_id AND i0.image_type_id = it0.image_type_id)
		JOIN image_type it1 ON(it1.abbr = 'rmi0')
		LEFT JOIN image i1 ON(i1.article_id = a.article_id AND i1.image_type_id = it1.image_type_id)
		ORDER BY a.created DESC`)

	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var picturePrefixes = [...]string{"rma0", "rmi0"}

	for rows.Next() {
		var a Article

		var created rawTime

		var pictures [2]sql.NullString

		err := rows.Scan(&a.Id, &a.Name, &created, &pictures[0], &pictures[1])
		if err != nil {
			panic(err.Error())
		}

		a.Created = created.Time()
		a.Pictures = make(PictureMap)

		for i, p := range pictures {
			if p.Valid {
				a.Pictures[picturePrefixes[i]] = Picture{
					Path: strings.Trim(p.String, "\n\r "),
					Type: picturePrefixes[i],
				}
			}
		}
		articles = append(articles, a)
	}

	return
}
