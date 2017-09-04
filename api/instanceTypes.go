package api

import (
	"gopkg.in/guregu/null.v3"
)

type Instance struct {
	InstanceId   int         `json:"instanceId"   db:"instance_id"`
	ArticleId    int         `json:"articleId"    db:"article_id"`
	LengthMm     null.Int    `json:"lengthMm"     db:"length_mm"`
	WidthMm      null.Int    `json:"widthMm"      db:"width_mm"`
	Picture0     null.String `json:"picture0"     db:"picture0"`
	PriceCchf    null.Int    `json:"priceCchf"    db:"price_cchf"`
	DiscountCchf null.Int    `json:"discountCchf" db:"discount_cchf"`
}

type InstanceColor struct {
	ColorId   int            `json:"colorId"      db:"color_id"`
	ColorName null.String    `json:"colorName"    db:"color_name"`
	CatName   null.String    `json:"catName"      db:"cat_name"`
	Hex       null.String    `json:"hex"          db:"hex"`
}

type InstanceImage struct {
	ImageType null.String    `json:"imageType"    db:"image_type"`
	Path      string         `json:"path"         db:"path"`
}

type InstanceDetailSpecifics struct {
	PatternId   null.Int        `json:"patternId"     db:"pattern_id"`
	ArticleName null.String     `json:"articleName"   db:"article_name"`
	ArticleDesc null.String     `json:"articleDesc"   db:"article_desc"`
	HeightMm    null.Int        `json:"heightMm"      db:"height_mm"`
	NumbStrings null.Int        `json:"numbStrings"   db:"numb_strings"`
	Colors      []InstanceColor `json:"colors"`
	GarnType    null.String     `json:"garnType"      db:"garn_type"`
	Images      []InstanceImage `json:"images"`
}

type InstanceDetail struct {
	Instance
	InstanceDetailSpecifics
}
