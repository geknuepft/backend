package articles

type FilterValue struct {
	FilterId int   `json:"filter_id"`
	Values   []int `json:"values"`
}

type FilterValues []FilterValue

func (filterValue FilterValue) GetFilter() (filter Filter, err error) {
	return GetFilterById(filterValue.FilterId)
}
