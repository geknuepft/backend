package api

import (
	"github.com/geknuepft/backend/database"
	"log"
	"errors"
)

func GetInstanceDetailById(instanceId int) (instanceDetail InstanceDetail, err error) {
	instance, err := GetInstanceById(instanceId)

	if err != nil {
		return
	}

	qs := getInstanceDetailQs(
		"i.instance_id = :instance_id",
		"",
	)

	return getInstanceDetailByQs(instance, qs, map[string]interface{}{"instance_id": instanceId})
}

// get query string to fetch instances from db
func getInstanceDetailQs(where, orderBy string) (qs string) {
	qs = `
    SELECT
      COALESCE(a.article_de, cat.category_de) article_name
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
    GROUP BY NULL
    ` + database.IfNotEmpty("ORDER BY ", orderBy)

	log.Printf("getInstanceDetailQs qs=%s", qs)

	return
}

type instanceDetailSpecificsRow struct {
	InstanceDetailSpecifics
}

func getInstanceDetailByQs(instance Instance, qs string, qArgs interface{}) (instanceDetail InstanceDetail, err error) {
	var db = database.GetDbX()
	defer db.Close()

	instanceDetail.Instance = instance

	rows, dbErr := db.NamedQuery(qs, qArgs)
	if dbErr != nil {
		log.Printf(dbErr.Error())
		return
	}

	var instanceDetailSpecificsRow = instanceDetailSpecificsRow{}

	if !rows.Next() {
		err = errors.New("Empty result")
	}

	dbErr = rows.StructScan(&instanceDetailSpecificsRow)
	if dbErr != nil {
		log.Printf(dbErr.Error())
		return
	}

	a := &instanceDetailSpecificsRow.InstanceDetailSpecifics

	instanceDetail.InstanceDetailSpecifics = *a

	return
}
