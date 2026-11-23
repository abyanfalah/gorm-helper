package gormhelper

import (
	"fmt"

	"gorm.io/gorm"
)

func Paginate[T any](queryChain *gorm.DB, req PaginationRequest) (PaginationWrapper[T], error) {
	var data []T
	var totalCount int64

	page := max(req.Page, 1)
	size := max(req.Size, 1)
	offset := (page - 1) * size

	// Count total records on a cloned session without limit/offset
	countDB := queryChain.Session(&gorm.Session{})
	if err := countDB.Count(&totalCount).Error; err != nil {
		return PaginationWrapper[T]{}, fmt.Errorf("count failed: %w", err)
	}

	if err := queryChain.Limit(size).Offset(offset).Find(&data).Error; err != nil {
		return PaginationWrapper[T]{}, fmt.Errorf("data fetch failed: %w", err)
	}

	totalPages := (totalCount + int64(size) - 1) / int64(size)

	return PaginationWrapper[T]{
		Total:     totalCount,
		TotalPage: totalPages,
		Page:      page,
		Showing:   len(data),
		Data:      data,
	}, nil
}

// PaginateRaw returns paginated data without the metadata wrapper.
func PaginateRaw[T any](queryChain *gorm.DB, req PaginationRequest) ([]T, error) {
	var data []T

	page := max(req.Page, 1)
	size := max(req.Size, 1)
	offset := (page - 1) * size

	if err := queryChain.Limit(size).Offset(offset).Find(&data).Error; err != nil {
		return nil, fmt.Errorf("data fetch failed: %w", err)
	}

	return data, nil
}