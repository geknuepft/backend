package main

import (
 "gopkg.in/guregu/null.v3"
    "time"
)

type Article struct {
    Id       int        `json:"article_id"`
    Name     string     `json:"name"`
    Created  time.Time  `json:"created"`
    Pictures PictureMap `json:"pictures"`
    InstanceId null.Int `json:"instance_id"`
}

type Articles []Article
