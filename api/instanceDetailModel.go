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
      a.article_id,
      p.pattern_id,
      a.article_desc_de                       article_desc,
      i.height_mm,
      p.numb_strings,
      (
        SELECT product_type_de
        FROM product_type     pt
        JOIN product_type_cat ptc ON(ptc.product_type_cat_id = pt.product_type_cat_id)
        JOIN product          p   ON(p.product_type_id = pt.product_type_id)
        JOIN material         m   ON(m.product_id = p.product_id)
        JOIN component        c   ON(c.material_id = m.material_id)
        WHERE ptc.abbr = 'G' AND c.article_id = a.article_id
        -- make query deterministic
        ORDER BY pt.product_type_id ASC
        LIMIT 1
      ) garn_type
    FROM instance          i
    JOIN article           a    ON(a.article_id = i.article_id)
    JOIN pattern           p    ON(p.pattern_id = a.pattern_id)
    JOIN category          cat  ON(cat.category_id = a.category_id)
    ` + database.IfNotEmpty("WHERE ", where);
	log.Printf("getInstanceDetailQs qs=%s", qs)

	return
}

func getInstanceDetailSpecifics(qs string, qArgs interface{}) (instanceDetailSpecifics InstanceDetailSpecifics, err error) {
	var db = database.GetDbX()

	// fetch instance details
	rows, dbErr := db.NamedQuery(qs, qArgs)
	defer rows.Close()
	if dbErr != nil {
		log.Printf(dbErr.Error())
		return
	}

	if !rows.Next() {
		err = errors.New("Empty result")
	}

	dbErr = rows.StructScan(&instanceDetailSpecifics)
	if dbErr != nil {
		log.Printf(dbErr.Error())
		return
	}

	return
}

func getColorsQs() (qs string) {

	qs = `
SELECT DISTINCT
  col.color_id               color_id,
  col.color_de               color_name,
  ccat.color_cat_de          cat_name,
  COALESCE(col.hex,ccat.hex) hex
FROM instance          i
JOIN article           a    ON(a.article_id = i.article_id)
JOIN component         co   ON(co.article_id = a.article_id)
JOIN material          m    ON(m.material_id = co.material_id)
JOIN material_x_color  mxc  ON(mxc.material_id = co.material_id)
JOIN color             col  ON(col.color_id = mxc.color_id)
JOIN color_cat         ccat ON(ccat.color_cat_id = col.color_cat_id)
WHERE i.instance_id = :instance_id
ORDER BY co.position ASC`

	return
}

func getColors(instanceId int) (instanceColors []InstanceColor, err error) {
	var db = database.GetDbX()

	// fetch instance details
	rows, dbErr := db.NamedQuery(getColorsQs(), map[string]interface{}{"instance_id": instanceId})
	defer rows.Close()

	if dbErr != nil {
		log.Printf(dbErr.Error())
		return
	}

	var instanceColor = InstanceColor{}

	for rows.Next() {
		err := rows.StructScan(&instanceColor)
		if err != nil {
			log.Printf(err.Error())
			break
		}

		instanceColors = append(instanceColors, instanceColor)
	}

	return
}

func getImagesQs() (qs string) {

	qs = `
SELECT
  it.image_type_de image_type,
  im.path
FROM instance          i
JOIN article           a    ON(a.article_id = i.article_id)
JOIN image             im   ON(im.article_id = a.article_id)
JOIN image_type        it   ON(it.image_type_id = im.image_type_id)
WHERE i.instance_id = :instance_id AND it.gallery
ORDER BY it.image_type_id ASC`

	return
}

func getImages(instanceId int) (instanceImages []InstanceImage, err error) {
	var db = database.GetDbX()

	// fetch instance details
	rows, dbErr := db.NamedQuery(getImagesQs(), map[string]interface{}{"instance_id": instanceId})
	defer rows.Close()

	if dbErr != nil {
		log.Printf(dbErr.Error())
		return
	}

	var instanceImage = InstanceImage{}

	for rows.Next() {
		err := rows.StructScan(&instanceImage)
		if err != nil {
			log.Printf(err.Error())
			break
		}

		instanceImages = append(instanceImages, instanceImage)
	}

	return
}

func getInstanceDetailByQs(instance Instance, qs string, qArgs interface{}) (instanceDetail InstanceDetail, err error) {
	// use general instance query data
	instanceDetail.Instance = instance

	// fetch instance details
	instanceDetail.InstanceDetailSpecifics, err = getInstanceDetailSpecifics(qs, qArgs);
	if (err != nil) {
		return
	}

	// fetch colors
	instanceDetail.Colors, err = getColors(instance.InstanceId)
	if (err != nil) {
		return
	}

	// fetch pictures
	instanceDetail.Images, err = getImages(instance.InstanceId)
	if (err != nil) {
		return
	}

	return
}
