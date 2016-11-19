package main

import (
	"gopkg.in/guregu/null.v3"
)

type Article struct {
	Id           int         `json:"article_id"`
	Name         string      `json:"name"`
	Pictures     PictureMap  `json:"pictures"`
	InstanceId   null.Int    `json:"instance_id"`
	LengthMm     null.Int    `json:"length_mm"`
	WidthMm      null.Int    `json:"width_mm"`
	HeightMm     null.Int    `json:"height_mm"`
	PriceCchf    null.Int    `json:"price_cchf"`
	CollectionDe null.String `json:"collection_de"`
}

type Articles []Article
