package old

import (
	"fmt"
	"gopkg.in/guregu/null.v3"
	"log"
	"strconv"
	"strings"
	"github.com/geknuepft/backend/database"
)

type Article struct {
	Id        int         `json:"article_id"    db:"article_id"`
	Name      string      `json:"name"          db:"article_name"`
	Picture   null.String `json:"pictures"      db:"path0"`
	WidthMm   null.Int    `json:"width_mm"      db:"width_mm"`
	PriceCchf null.Int    `json:"price_cchf"    db:"price_cchf"`
}

type ArticleRow struct {
	Article
}

type Articles []Article

func GetArticlesByFilterValues(filterValues FilterValues, lengthMm int) (articles Articles) {
	wheres, qArgs := getArticleWhere(filterValues)
	where := strings.Join(wheres, " AND ")

	return getArticlesByQs(
		getArticleQs(
			where,
			"a.created DESC, a.article_id ASC",
			lengthMm,
		),
		qArgs,
	)
}

func GetArticles() (articles Articles) {
	qs := getArticleQs(
		"",
		"a.created DESC",
		180,
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

		articles = append(articles, *a)
	}

	return
}

func getArticleQs(where, orderBy string, lengthMm int) (qs string) {
	lengthMmStr := strconv.Itoa(lengthMm)

	qs = `
        SELECT
        -- article fields
        a.article_id,
        COALESCE(a.article_name_de, cat.category_de) article_name,
        i0.path path0,
        -- instance fields
        ROUND(p.width_mm * AVG(pr.pattern_factor_width)) width_mm,
        ROUND(
          (
            (p.price_cchf + ` + lengthMmStr + ` * p.price_cchf_cm / 10) +
            COALESCE(
              FLOOR(p.numb_pearls + ` + lengthMmStr + ` * p.numb_pearls_10cm / 100)
              * MAX(m.price_pp_cchf),
              0
            )
          ),-1
        ) price_cchf
        FROM article           a
        JOIN category          cat  ON(cat.category_id = a.category_id)
        JOIN image_type        it0  ON(it0.image_type = 'square-narrow_white_single_plan_setup0')
        JOIN image             i0   ON(i0.article_id = a.article_id AND i0.image_type_id = it0.image_type_id)
        JOIN collection        ac   ON(ac.collection_id = a.collection_id)
        JOIN component         co   ON(co.article_id = a.article_id)
        JOIN material          m    ON(m.material_id = co.material_id)
        JOIN product           pr   ON(pr.product_id = m.product_id)
        JOIN material_x_color  mxc  ON(mxc.material_id = co.material_id)
        JOIN color             col  ON(col.color_id = mxc.color_id)
        JOIN color_cat         ccat ON(ccat.color_cat_id = col.color_cat_id)
        JOIN pattern           p    ON(p.pattern_id = a.pattern_id)
        ` + database.IfNotEmpty("WHERE ", where) + `
        GROUP BY a.article_id 
        ` + database.IfNotEmpty("ORDER BY ", orderBy)

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
	case "range":
		qValues := make([]string, 0, len(filterValue.Values))

		fieldName := fmt.Sprintf("%s.%s", alias, filter.DbColumn)

		for _, value := range filterValue.Values {
			if fr, ok := filter.Range[strconv.Itoa(value)]; ok {
				qValues = append(qValues, fr.GetSql(fieldName))
			}
		}
		where = "(" + strings.Join(qValues, " OR ") + ")"
	default:
		panic("article: getArticleJoin: unsupported type: " + filter.FilterType)
	}

	return
}
