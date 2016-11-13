package main

import (
    "database/sql"
    "time"
)

type Article struct {
    Id       int        `json:"article_id"`
    Name     string     `json:"name"`
    Created  time.Time  `json:"created"`
    Pictures PictureMap `json:"pictures"`
    InstanceId sql.NullInt64 `json:"instance_id"`
}

type Articles []Article
