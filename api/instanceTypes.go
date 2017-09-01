package api

import (
	"gopkg.in/guregu/null.v3"
)

type Instance struct {
	InstanceId   int         `json:"instance_id"   db:"instance_id"`
	ArticleId    int         `json:"article_id"    db:"article_id"`
	LengthMm     null.Int    `json:"length_mm"     db:"length_mm"`
	WidthMm      null.Int    `json:"width_mm"      db:"width_mm"`
	Picture0     null.String `json:"picture0"      db:"picture0"`
	PriceCchf    null.Int    `json:"price_cchf"    db:"price_cchf"`
	DiscountCchf null.Int    `json:"discount_cchf" db:"discount_cchf"`
}

type InstanceDetailSpecifics struct {
	ArticleName null.String `json:"article_name"   db:"article_name"`
}

type InstanceDetail struct {
	Instance
	InstanceDetailSpecifics
}
