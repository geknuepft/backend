package old

import (
	"errors"
	"fmt"
	"gopkg.in/guregu/null.v3"
	"log"
	"strconv"
	"github.com/geknuepft/backend/sql"
)

type Filter struct {
	Id         int64                     `json:"filter_id"   db:"filter_id"`
	FilterType string                    `json:"filter_type" db:"filter_type"`
	Filter     string                    `json:"filter"      db:"filter"`
	FilterDe   string                    `json:"filter_de"   db:"filter_de"`
	DbTable    string                    `json:"-"           db:"db_table"`
	DbColumn   string                    `json:"-"           db:"db_column"`
	Category   map[string]FilterCategory `json:"category"`
	Range      map[string]FilterRange    `json:"range"`
}

type FilterRange struct {
	FilterRangeId null.Int `json:"filter_range_id" db:"filter_range_id"`
	RangeGeq      null.Int `json:"geq"             db:"range_geq"`
	RangeLeq      null.Int `json:"leq"             db:"range_leq"`
}

type FilterCategory struct {
	FilterCategoryId null.Int    `json:"filter_category_id" db:"filter_category_id"`
	CategoryDe       null.String `json:"de"                 db:"filter_category_de"`
}

type FilterRow struct {
	Filter
	FilterRange
	FilterCategory
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
			filter.Category = make(map[string]FilterCategory, 10)
		}

		// append a range entry when filter_range_id IS NOT NULL
		if filterRow.FilterRangeId.Valid {
			filter.Range[strconv.FormatInt(filterRow.FilterRangeId.Int64, 10)] = filterRow.FilterRange
		}

		// append a category entry when filter_category_id IS NOT NULL
		if filterRow.FilterCategoryId.Valid {
			filter.Category[strconv.FormatInt(filterRow.FilterCategoryId.Int64, 10)] = filterRow.FilterCategory
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
			"filter_id ASC, filter_range_id ASC, filter_category_id ASC",
		),
		map[string]interface{}{},
	)
}

func getFilterQs(where, orderBy string) (qs string) {
	qs = "SELECT filter_id, filter_type, filter, filter_de, " +
		"db_table, db_column, " +
		"filter_range_id, range_geq, range_leq, " +
	// filter_category fields
		"COALESCE(category.category_id, color_cat.color_cat_id) filter_category_id, " +
		"COALESCE(category.category_de, color_cat.color_cat_de) filter_category_de " +
		"FROM filter f " +
		"JOIN filter_type USING(filter_type_id) " +
		"LEFT JOIN filter_range USING(filter_id) " +
		"LEFT JOIN category ON(f.db_table = 'category') " +
		"LEFT JOIN color_cat ON(f.db_table = 'color_cat') " +
		sql.IfNotEmpty("WHERE ", where) + " " +
		sql.IfNotEmpty("ORDER BY ", orderBy)

	log.Printf("qs=%s", qs)

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
