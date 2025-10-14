package utils

import "math"

type Pagination struct {
	Page        int32 `json:page`
	Limit       int32 `json:limit`
	TotalRecord int32 `json:total_records`
	TotalPage   int32 `json:total_pages`
	HasNext     bool  `json:has_next`
	HasPrev     bool  `json:has_prev`
}

func NewPagination(page, limit, totalRecords int32) *Pagination {

	// Tổng số trang = ceil(tổng số record / số phần tử trên 1 trang)
	totalPages := math.Ceil(float64(totalRecords) / float64(limit))

	return &Pagination{
		Page:        page,
		Limit:       limit,
		TotalRecord: totalRecords,
		TotalPage:   totalPages,
		HasNext:     page < int32(totalPages),
		HasPrev:     page > 1,
	}
}
