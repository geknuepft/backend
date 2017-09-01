package api

import (
	"github.com/geknuepft/backend/database"
	"log"
	"strings"
	"errors"
)

func GetInstances() (instances []Instance) {
	wheres, qArgs := getInstanceWhere()
	where := strings.Join(wheres, " AND ")

	return getInstancesByQs(
		getInstanceQs(
			where,
			"created DESC, instance_id ASC",
		),
		qArgs,
	)
}

func GetInstanceById(instanceId int) (instance Instance, err error) {
	instances := getInstancesByQs(
		getInstanceQs(
			"i.instance_id = :instance_id",
			"",
		),
		map[string]interface{}{"instance_id": instanceId},
	)

	if len(instances) < 1 {
		err = errors.New("not found")
		return
	}

	instance = instances[0]
	return
}

type instanceRow struct {
	Instance
}

func getInstancesByQs(qs string, qArgs interface{}) (instances []Instance) {
	var db = database.GetDbX()
	defer db.Close()

	rows, err := db.NamedQuery(qs, qArgs)
	if err != nil {
		log.Printf(err.Error())
	}

	var instanceRow = instanceRow{}

	for rows.Next() {
		err := rows.StructScan(&instanceRow)
		if err != nil {
			log.Printf(err.Error())
			break
		}

		a := &instanceRow.Instance

		instances = append(instances, *a)
	}

	return
}

// get query string to fetch instances from db
func getInstanceQs(where, orderBy string) (qs string) {
	qs = `
SELECT
  instance_id,
  article_id,
  length_mm,
  width_mm,
  picture0,
  price_cchf,
  calc_instance_discount_cchf(price_cchf) discount_cchf
FROM (
  SELECT
    f0.*,
    calc_price_cchf(
      f0.pattern_id,
      f0.length_mm,
      f0.width_mm,
      f0.price_pp_cchf
    ) price_cchf
  FROM (
    SELECT
      i.instance_id,
      i.article_id,
      i.length_mm,  -- go logs warning if this is null
      COALESCE(
        i.width_mm,
        ROUND(p.width_mm * AVG(pr.pattern_factor))
      ) width_mm,    -- go logs warning if this is null
      AVG(m.price_pp_cchf) price_pp_cchf,
      i0.path picture0,
      i.created,
      p.pattern_id
    FROM instance          i
    JOIN article           a    ON(a.article_id = i.article_id)
    JOIN pattern           p    ON(p.pattern_id = a.pattern_id)
    JOIN category          cat  ON(cat.category_id = a.category_id)
    JOIN image_type        it0  ON(it0.abbr = 'cma0')
    JOIN image             i0   ON(i0.article_id = a.article_id AND i0.image_type_id = it0.image_type_id)
    JOIN collection        ac   ON(ac.collection_id = a.collection_id)
    JOIN component         co   ON(co.article_id = a.article_id)
    JOIN material          m    ON(m.material_id = co.material_id)
    JOIN product           pr   ON(pr.product_id = m.product_id)
    JOIN material_x_color  mxc  ON(mxc.material_id = co.material_id)
    JOIN color             col  ON(col.color_id = mxc.color_id)
    JOIN color_cat         ccat ON(ccat.color_cat_id = col.color_cat_id)
    ` + database.IfNotEmpty("WHERE ", where) + `
    GROUP BY i.instance_id
    LIMIT 2048
  ) f0
) f1
` + database.IfNotEmpty("ORDER BY ", orderBy)

	log.Printf("getInstanceQs qs=%s", qs)

	return
}

func getInstanceWhere() (wheres []string, qArgs map[string]interface{}) {
	wheres = make([]string, 0)

	wheres = append(wheres, "cat.abbr IN('A','S','B','H','T','E','AP','F','FP')")
	wheres = append(wheres, "i.length_mm IS NOT NULL")

	return
}
