package gormhelper

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/* this function only saves like... a line. usage is not obligatory */
func FindByUUID[T any](ctx context.Context, db *gorm.DB, id uuid.UUID) (*T, error) {
	if id == uuid.Nil {
		return nil, nil
	}

	var entity T
	err := db.WithContext(ctx).First(&entity, id).Error
	return &entity, err
}

func SearchWithPagination[T any](ctx context.Context, db *gorm.DB, req *SearchPaginationRequest, fieldsToSearch ...string) (PaginationWrapper[T], error) {
	query := db.WithContext(ctx).Model(new(T)).Scopes(
		SearchScope(req.SearchRequest, fieldsToSearch...),
	)

	result, err := Paginate[T](query, req.PaginationRequest)
	if err != nil {
		return PaginationWrapper[T]{}, err
	}

	return result, nil
}

func FilteredSearchWithPagination[T any](ctx context.Context, db *gorm.DB, req *SearchPaginationRequest, fieldFilters map[string]any, fieldsToSearch ...string) (PaginationWrapper[T], error) {
	query := db.WithContext(ctx).Model(new(T)).
		Scopes(FieldFilterScope(fieldFilters)).
		Scopes(SearchScope(req.SearchRequest, fieldsToSearch...))

	return Paginate[T](query, req.PaginationRequest)

}
func ModelFilteredSearchWithPagination[T any](
	ctx context.Context,
	db *gorm.DB,
	req *SearchPaginationRequest,
	filter T,
	fieldsToSearch ...string,
) (PaginationWrapper[T], error) {
	query := db.
		WithContext(ctx).
		Model(new(T)).
		Where(&filter).
		Scopes(SearchScope(req.SearchRequest, fieldsToSearch...))

	return Paginate[T](query, req.PaginationRequest)
}
func IsExistsBy[T any](ctx context.Context, db *gorm.DB, field string, value any) (bool, error) {
	var count int64
	err := db.WithContext(ctx).Model(new(T)).Where(fmt.Sprintf("%s = ?", field), value).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func IsExistsById[T any](ctx context.Context, db *gorm.DB, id uint) (bool, error) {
	var count int64
	err := db.WithContext(ctx).Model(new(T)).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, err
}

func FindBy[T any](ctx context.Context, db *gorm.DB, field string, value any) (*T, error) {
	model := new(T)
	err := db.WithContext(ctx).Model(model).Where(fmt.Sprintf("%s = ?", field), value).Find(&model).Error

	return model, err
}
