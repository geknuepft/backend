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

type InstanceDetailSpecifics struct {
	ArticleName null.String `json:"articleName"   db:"article_name"`
}

type InstanceDetail struct {
	Instance
	InstanceDetailSpecifics
}
