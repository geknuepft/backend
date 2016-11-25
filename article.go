package main

import (
	"gopkg.in/guregu/null.v3"
	"log"
	"strings"
)

type Article struct {
	Id           int         `json:"article_id"    db:"article_id"`
	Name         string      `json:"name"          db:"article_name"`
	Pictures     PictureMap  `json:"pictures"`
	InstanceId   null.Int    `json:"instance_id"   db:"instance_id"`
	LengthMm     null.Int    `json:"length_mm"     db:"length_mm""`
	WidthMm      null.Int    `json:"width_mm"      db:"width_mm"`
	HeightMm     null.Int    `json:"height_mm"     db:"height_mm"`
	PriceCchf    null.Int    `json:"price_cchf"    db:"price_cchf"`
	CollectionDe null.String `json:"collection_de" db:"collection_de"`
}

type ArticleRow struct {
	Article
	Path0 null.String `db:"path0"`
}

type Articles []Article

var picturePrefixes = [...]string{"cma0"}

func GetArticlesByFilterValues(filterValues FilterValues) (articles Articles) {

	where, qArgs := getArticleWhere(filterValues)

	return getArticlesByQs(
		getArticleQs(
			where,
			"a.created DESC",
		),
		qArgs,
	)
}

func GetArticles() (articles Articles) {

	qs := getArticleQs(
		"",
		"a.created DESC",
	)

	return getArticlesByQs(qs, map[string]interface{}{})
}

func getArticlesByQs(qs string, qArgs interface{}) (articles Articles) {

	var db = getDbX()
	defer db.Close()

	rows, err := db.NamedQuery(qs, qArgs)
	if err != nil {
		log.Printf(err.Error())
	}

	var articleRow = ArticleRow{}

	for rows.Next() {
		err := rows.StructScan(&articleRow)
		if err != nil {
			log.Printf(err.Error())
			break
		}

		a := &articleRow.Article

		if articleRow.Path0.Valid {
			a.Pictures = make(PictureMap)
			a.Pictures[picturePrefixes[0]] = Picture{
				Path: strings.Trim(articleRow.Path0.String, "\n\r "),
				Type: picturePrefixes[0],
			}
		}

		articles = append(articles, *a)
	}

	return
}

func getArticleQs(where, orderBy string) (qs string) {
	qs = `SELECT
        -- article fields
        a.article_id,
        COALESCE(a.article_name_de, c.category_de) article_name,
        i0.path path0,
        -- instance fields
        i.instance_id,
        i.length_mm,
        i.width_mm,
        i.height_mm,
        i.price_cchf,
        ic.collection_de
        FROM article a
        JOIN category c ON(c.category_id = a.category_id)
        JOIN image_type it0 ON(it0.abbr = '` + picturePrefixes[0] + `')
        JOIN image i0 ON(i0.article_id = a.article_id AND i0.image_type_id = it0.image_type_id)
        LEFT JOIN instance i ON(i.article_id = a.article_id)
        LEFT JOIN collection ic ON(ic.collection_id = i.collection_id)` +
		IfNotEmpty("WHERE ", where) + `
        GROUP BY a.article_id ` + // in case an article has >1 cma0
		IfNotEmpty("ORDER BY ", orderBy)

	return
}

func getArticleWhere(filterValues FilterValues) (where string, qArgs map[string]interface{}) {

	qArgs = make(map[string]interface{}, len(filterValues))

	for _, filterValue := range filterValues {
		log.Printf("%v", filterValue)

	}

	return
}
