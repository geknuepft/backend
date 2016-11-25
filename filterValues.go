package main

type FilterValue struct {
	FilterId int   `json:"filter_id"`
	Values   []int `json:"values"`
}

type FilterValues []FilterValue
