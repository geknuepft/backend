package main

import (
	"gopkg.in/guregu/null.v3"
	"log"
)

type Filter struct {
	Id         int           `json:"filter_id"   db:"filter_id"`
	FilterType string        `json:"filter_type" db:"filter_type"`
	Filter     string        `json:"filter"      db:"filter"`
	FilterDe   string        `json:"filter_de"   db:"filter_de"`
	Range      []FilterRange `json:"filer_range"`
}

type FilterRange struct {
	FilterRangeId null.Int `db:"filter_range_id"`
	RangeGeq      null.Int `json:"range_geq" db:"range_geq"`
	RangeLeq      null.Int `json:"range_leq" db:"range_leq"`
}

type FilterRow struct {
	Filter
	FilterRange
}

type Filters []Filter

func getFilterByQs(qs string, qArgs ...interface{}) (filters Filters) {
	var db = getDbX()
	defer db.Close()

	filterRow := FilterRow{}

	rows, err := db.Queryx(qs, qArgs...)
	if err != nil {
		log.Printf(err.Error())
	}

	var filter *Filter
	for rows.Next() {
		err := rows.StructScan(&filterRow)
		if err != nil {
			log.Printf(err.Error())
			break
		}

		// create a new filter whenever filter_id changes
		if len(filters) < 1 || filter.Id != filterRow.Id {
			filters = append(filters, filterRow.Filter)
			filter = &filters[len(filters)-1]
		}

		// append a range entry when filter_range_id IS NOT NULL
		if filterRow.FilterRangeId.Valid {
			filter.Range = append(filter.Range, filterRow.FilterRange)
		}
	}

	return
}

func GetFilterById(id int) (filters Filters) {
	return getFilterByQs(
		getQs(
			"filter_id = ?",
			"",
		),
		id,
	)
}

func GetFilters() (filters Filters) {
	return getFilterByQs(
		getQs(
			"",
			"filter_id ASC, filter_range_id ASC",
		),
	)
}

func getQs(where, orderBy string) (qs string) {
	qs = "SELECT filter_id, filter_type, filter, filter_de, " +
		"filter_range_id, range_geq, range_leq " +
		"FROM filter  " +
		"JOIN filter_type USING(filter_type_id) " +
		"LEFT JOIN filter_range USING(filter_id) " +
		ifNotEmpty("WHERE ", where) + " " +
		ifNotEmpty("ORDER BY ", orderBy)
	return
}

func ifNotEmpty(prefix, value string) string {
	if len(value) > 0 {
		return prefix + value
	}
	return ""
}
