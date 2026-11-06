package gormhelper

type ResponseWrapper struct {
}

type PaginationWrapper[T any] struct {
	Total     int64 `json:"total"`
	TotalPage int64 `json:"totalPage"`
	Page      int   `json:"page"`
	Showing   int   `json:"showing"`
	Data      []T   `json:"data"`
}
