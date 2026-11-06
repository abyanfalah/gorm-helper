package github.com/abyanfalah/gorm-helper

import (
	"fmt"

	"gorm.io/gorm"
)

func Paginate[T any](queryChain *gorm.DB, req PaginationRequest) (PaginationWrapper[T], error) {
	var data []T
	var totalCount int64

	zeroBasedPage := req.Page - 1

	// Count total records
	countDB := queryChain.Session(&gorm.Session{}) // clone the query without limits
	if err := countDB.Count(&totalCount).Error; err != nil {
		return PaginationWrapper[T]{}, fmt.Errorf("count failed: %w", err)
	}

	// Apply limit and zero based offset for pagination
	offset := (zeroBasedPage * req.Size)

	if err := queryChain.Limit(req.Size).Offset(offset).Find(&data).Error; err != nil {
		return PaginationWrapper[T]{}, fmt.Errorf("data fetch failed: %w", err)
	}

	totalPages := int((totalCount + int64(req.Size) - 1) / int64(req.Size))

	return PaginationWrapper[T]{
		Total:     totalCount,
		TotalPage: int64(totalPages),
		Page:      req.Page,
		Showing:   len(data),
		Data:      data,
	}, nil
}

/* get paginated data without wrapper */
func PaginateRaw[T any](queryChain *gorm.DB, req PaginationRequest) ([]T, error) {
	var data []T
	zeroBasedPage := req.Page - 1
	offset := (zeroBasedPage * req.Size)

	if err := queryChain.Limit(req.Size).Offset(offset).Find(&data).Error; err != nil {
		return nil, fmt.Errorf("data fetch failed: %w", err)
	}

	return data, nil
}
