package utils

type Pagination struct {
	Page        int32 `json:page`
	Limit       int32 `json:limit`
	TotalRecord int32 `json:total_records`
	TotalPage   int32 `json:total_pages`
	HasNext     bool  `json:has_next`
	HasPrev     bool  `json:has_prev`
}
