package api

import (
	"gopkg.in/guregu/null.v3"
)

type Instance struct {
	InstanceId   int         `json:"instance_id"   db:"instance_id"`
	LengthMm     null.Int    `json:"length_mm"     db:"length_mm"`
	WidthMm      null.Int    `json:"width_mm"      db:"width_mm"`
	Picture      null.String `json:"pictures"      db:"path0"`
	PriceCchf    null.Int    `json:"price_cchf"    db:"price_cchf"`
	DiscountCchf null.Int    `json:"discount_cchf" db:"discount_cchf"`
}
