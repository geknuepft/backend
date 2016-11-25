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
	RangeGeq null.Int `json:"range_geq" db:"range_geq"`
	RangeLeq null.Int `json:"range_leq" db:"range_leq"`
}

type FilterRow struct {
	Filter
	FilterRangeId null.Int `db:"filter_range_id"`
	FilterRange
}

type Filters []Filter

func GetFilters() (filters Filters) {

	var db = getDbX()
	defer db.Close()

	var qs = "SELECT filter_id, filter_type, filter, filter_de, " +
		"filter_range_id, range_geq, range_leq " +
		"FROM filter  " +
		"JOIN filter_type USING(filter_type_id) " +
		"LEFT JOIN filter_range USING(filter_id) " +
		"ORDER BY filter_id ASC, filter_range_id ASC"

	filterRow := FilterRow{}

	rows, err := db.Queryx(qs)
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
