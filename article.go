package main

import (
	"fmt"
	"gopkg.in/guregu/null.v3"
	"log"
	"strconv"
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
	wheres, qArgs := getArticleWhere(filterValues)
	where := strings.Join(wheres, " AND ")

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
	qs = `
        SELECT
        -- article fields
        a.article_id,
        COALESCE(a.article_name_de, cat.category_de) article_name,
        i0.path path0,
        -- instance fields
        i.instance_id,
        i.length_mm,
        i.width_mm,
        i.height_mm,
        i.price_cchf,
        ic.collection_de
        FROM article a
        JOIN category          cat  ON(cat.category_id = a.category_id)
        JOIN image_type        it0  ON(it0.abbr = '` + picturePrefixes[0] + `')
        JOIN image             i0   ON(i0.article_id = a.article_id AND i0.image_type_id = it0.image_type_id)
        JOIN instance          i    ON(i.article_id = a.article_id)
        JOIN collection        ic   ON(ic.collection_id = i.collection_id)
        JOIN component         co   ON(co.article_id = a.article_id)
        JOIN material_x_color  mxc  ON(mxc.material_id = co.material_id)
        JOIN color             col  ON(col.color_id = mxc.color_id)
        JOIN color_cat         ccat ON(ccat.color_cat_id = col.color_cat_id)
        ` + IfNotEmpty("WHERE ", where) + `
        GROUP BY a.article_id 
        ` + IfNotEmpty("ORDER BY ", orderBy)

	log.Printf("qs=%s", qs)

	return
}

func getArticleWhere(filterValues FilterValues) (wheres []string, qArgs map[string]interface{}) {
	wheres = make([]string, 0, len(filterValues))

	for _, filterValue := range filterValues {
		filter, err := filterValue.GetFilter()
		if err != nil {
			continue
		}
		where := getArticleWhereByFilter(filter, filterValue)
		wheres = append(wheres, where)
	}

	return
}

func getArticleWhereByFilter(filter Filter, filterValue FilterValue) (where string) {

	// check that DbTable is supported
	var alias string
	switch filter.DbTable {
	case "instance":
		alias = "i"
	case "category":
		alias = "cat"
	case "color":
		alias = "col"
	case "color_cat":
		alias = "ccat"
	default:
		panic("article: getArticleJoin: unknown table: " + filter.DbTable)
	}

	// check that the FilterType is supported
	switch filter.FilterType {
	case "id":
		qValues := make([]string, len(filterValue.Values))
		for i, value := range filterValue.Values {
			qValues[i] = strconv.Itoa(value)
		}
		in := strings.Join(qValues, ",")
		where = fmt.Sprintf("%s.%s IN(%s)", alias, filter.DbColumn, in)
	default:
		panic("article: getArticleJoin: unsupported type: " + filter.FilterType)
	}

	return
}
