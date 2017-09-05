package api

import (
	"github.com/geknuepft/backend/database"
	"log"
)

// get query string to fetch instances from db
func getColorCatQs() (qs string) {
	qs = `
  SELECT
      color_cat_id,
      color_cat_de name,
      hex
  FROM color_cat
  ORDER BY sort ASC`

	return
}

func getColorCats() (colorCats []ColorCat, err error) {
	var db = database.GetDbX()
	defer db.Close()

	// fetch instance details
	rows, dbErr := db.NamedQuery(getColorCatQs(), map[string]interface{}{})
	if dbErr != nil {
		log.Printf(dbErr.Error())
		return
	}

	var colorCat = ColorCat{}

	for rows.Next() {
		err := rows.StructScan(&colorCat)
		if err != nil {
			log.Printf(err.Error())
			break
		}

		colorCats = append(colorCats, colorCat)
	}

	return
}
