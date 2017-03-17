package main

import (
	"errors"
	"gopkg.in/guregu/null.v3"
	"log"
	"strings"
)

type ArticleDetail struct {
	Article
}

type ArticleDetailRow struct {
	ArticleDetail
	Path0 null.String `db:"path0"`
}

//var picturePrefixes = [...]string{"cma0"}

func GetArticleDetailById(articleId int) (articleDetail ArticleDetail, err error) {
	qs := getArticleDetailQs(
		"a.article_id = :article_id",
		"",
	)

	return getArticleDetailByQs(qs, map[string]interface{}{"article_id": articleId})
}

func getArticleDetailByQs(qs string, qArgs interface{}) (articleDetail ArticleDetail, err error) {

	var db = getDbX()
	defer db.Close()

	rows, dbErr := db.NamedQuery(qs, qArgs)
	if dbErr != nil {
		log.Printf(dbErr.Error())
		return
	}

	var articleDetailRow = ArticleDetailRow{}

	if !rows.Next() {
		err = errors.New("Empty result")
	}

	dbErr = rows.StructScan(&articleDetailRow)
	if dbErr != nil {
		log.Printf(dbErr.Error())
		return
	}

	a := &articleDetailRow.ArticleDetail

	if articleDetailRow.Path0.Valid {
		a.Pictures = make(PictureMap)
		a.Pictures[picturePrefixes[0]] = Picture{
			Path: strings.Trim(articleDetailRow.Path0.String, "\n\r "),
			Type: picturePrefixes[0],
		}
	}

	articleDetail = *a

	log.Printf("articleDetail=%v", articleDetail)

	return
}

func getArticleDetailQs(where, orderBy string) (qs string) {
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
        LEFT JOIN collection        ic   ON(ic.collection_id = a.collection_id)
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
