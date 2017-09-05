package api

type ColorCat struct {
	Name string `json:"name" db:"name"`
	Hex  string `json:"hex"  db:"hex"`
}
