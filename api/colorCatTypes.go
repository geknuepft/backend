package api

type ColorCat struct {
	ColorCatId int    `json:"colorCatId" db:"color_cat_id"`
	Name       string `json:"name"       db:"name"`
	Hex        string `json:"hex"        db:"hex"`
}
