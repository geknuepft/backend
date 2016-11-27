package main

import (
	"errors"
	"fmt"
	"gopkg.in/guregu/null.v3"
	"log"
	"strconv"
)

type Filter struct {
	Id         int64                  `json:"filter_id"   db:"filter_id"`
	FilterType string                 `json:"filter_type" db:"filter_type"`
	Filter     string                 `json:"filter"      db:"filter"`
	FilterDe   string                 `json:"filter_de"   db:"filter_de"`
	DbTable    string                 `json:"-"           db:"db_table"`
	DbColumn   string                 `json:"-"           db:"db_column"`
	Range      map[string]FilterRange `json:"filer_range"`
}

type FilterRange struct {
	FilterRangeId null.Int `json:"filter_range_id" db:"filter_range_id"`
	RangeGeq      null.Int `json:"range_geq"       db:"range_geq"`
	RangeLeq      null.Int `json:"range_leq"       db:"range_leq"`
}

type FilterRow struct {
	Filter
	FilterRange
}

type Filters []Filter

func getFilterByQs(qs string, qArgs interface{}) (filters Filters) {
	var db = getDbX()
	defer db.Close()

	filterRow := FilterRow{}

	rows, err := db.NamedQuery(qs, qArgs)
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
			filter.Range = make(map[string]FilterRange, 10)
		}

		// append a range entry when filter_range_id IS NOT NULL
		if filterRow.FilterRangeId.Valid {
			filter.Range[strconv.FormatInt(filterRow.FilterRangeId.Int64, 10)] = filterRow.FilterRange
		}
	}

	return
}

func GetFilterById(filterId int) (filter Filter, err error) {
	filters := getFilterByQs(
		getFilterQs(
			"filter_id = :filter_id",
			"",
		),
		map[string]interface{}{"filter_id": filterId},
	)

	if len(filters) < 1 {
		err = errors.New("Empty result")
	} else {
		filter = filters[0]
	}
	return
}

func GetFilters() (filters Filters) {
	return getFilterByQs(
		getFilterQs(
			"",
			"filter_id ASC, filter_range_id ASC",
		),
		map[string]interface{}{},
	)
}

func getFilterQs(where, orderBy string) (qs string) {
	qs = "SELECT filter_id, filter_type, filter, filter_de, " +
		"db_table, db_column, " +
		"filter_range_id, range_geq, range_leq " +
		"FROM filter  " +
		"JOIN filter_type USING(filter_type_id) " +
		"LEFT JOIN filter_range USING(filter_id) " +
		IfNotEmpty("WHERE ", where) + " " +
		IfNotEmpty("ORDER BY ", orderBy)
	return
}

func (fr FilterRange) GetSql(fieldName string) string {
	if fr.RangeGeq.Valid && fr.RangeLeq.Valid {
		return fmt.Sprintf("%s BETWEEN %d AND %d", fieldName, fr.RangeGeq.Int64, fr.RangeLeq.Int64)
	} else if fr.RangeGeq.Valid {
		return fmt.Sprintf("%s >= %d", fieldName, fr.RangeGeq.Int64)
	} else if fr.RangeGeq.Valid {
		return fmt.Sprintf("%s <= %d", fieldName, fr.RangeLeq.Int64)
	} else {
		return ""
	}
}
