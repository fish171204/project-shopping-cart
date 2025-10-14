package utils

import (
	"strconv"
)

type Pagination struct {
	Page        int32 `json:"page"`
	Limit       int32 `json:"limit"`
	TotalRecord int32 `json:"total_records"`
	TotalPage   int32 `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

func NewPagination(page, limit, totalRecords int32) *Pagination {
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		envLimit := GetEnv("LIMIT_ITEM_ON_PER_PAGE", "10")
		limitInt, err := strconv.Atoi(envLimit)
		if err != nil || limitInt <= 0 {
			limit = 10
		}
		limit = int32(limitInt)
	}

	// Tổng số trang = ceil(tổng số record / số phần tử trên 1 trang)
	// C1: totalPages := math.Ceil(float64(totalRecords) / float64(limit))
	// C2: int
	totalPages := (totalRecords + limit - 1) / limit

	return &Pagination{
		Page:        page,
		Limit:       limit,
		TotalRecord: totalRecords,
		TotalPage:   totalPages,
		HasNext:     page < totalPages,
		HasPrev:     page > 1,
	}
}
