package main

import (
	"time"
)

type Article struct {
	Id      int       `json:"article_id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

type Articles []Article
