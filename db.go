package main

import (
    "database/sql"
    "strings"
    "log"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

func getDb() *sql.DB {
    var err error
    var Db *sql.DB

    Db, err = sql.Open("mysql", "geknuepft:Er3cof4iesho@tcp(mysql-server:3306)/geknuepft")
    if err != nil {
        panic(err.Error())
    }

    err = Db.Ping()
    if err != nil {
        panic(err.Error())
    }

    return Db
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
    var Db *sql.DB
    Db = getDb()
    defer Db.Close()

    var qs string

    qs = `SELECT
        -- article fields
        a.article_id,
        COALESCE(a.article_name_de, '') article_name,
        a.created,
        i0.path p0
        -- instance fields
        -- i.instance_id,
        -- i.length_mm,
        -- i.width_mm,
        -- i.height_mm,
        -- i.price_cchf,
        -- ic.collection_de
        FROM article a
        JOIN image_type it0 ON(it0.abbr = 'cma0')
        JOIN image i0 ON(i0.article_id = a.article_id AND i0.image_type_id = it0.image_type_id)
        JOIN instance i ON(i.article_id = a.article_id)
        JOIN collection ic ON(ic.collection_id = i.collection_id)
        GROUP BY a.article_id -- in case an article has >1 cma0
        ORDER BY a.created DESC`

    log.Printf("%s", qs)

    rows, err := Db.Query(qs)

    if err != nil {
        panic(err.Error())
    }
    defer rows.Close()

    var picturePrefixes = [...]string{"cma0"}

    for rows.Next() {
        var a Article

        var created rawTime

        var pictures [1]sql.NullString

        err := rows.Scan(&a.Id, &a.Name, &created, &pictures[0])
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
