package gormhelper

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

func PaginationScope(req PaginationRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := max(req.Page, 1)
		size := max(req.Size, 1)
		offset := (page - 1) * size
		return db.Limit(size).Offset(offset)
	}
}

func ActiveStatusScope(req IsAktifRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if req.IsAktif == nil {
			return db
		}
		return db.Where("is_active = ?", *req.IsAktif)
	}
}

func IsDeletedScope(req IsDeleted) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if req.IsDeleted == nil {
			return db.Unscoped()
		}
		if *req.IsDeleted {
			return db.Unscoped().Where("deleted_at IS NOT NULL")
		}
		return db
	}
}

func SearchScope(req SearchRequest, fieldsToSearch ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if req.Search == "" || len(fieldsToSearch) == 0 {
			return db
		}

		pattern := "%" + req.Search + "%"
		conditions := make([]string, 0, len(fieldsToSearch))
		args := make([]any, 0, len(fieldsToSearch))

		for _, field := range fieldsToSearch {
			conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", field))
			args = append(args, pattern)
		}

		return db.Where(strings.Join(conditions, " OR "), args...)
	}
}

func FieldFilterScope(fieldFilters map[string]any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for field, filterValue := range fieldFilters {
			val := reflect.ValueOf(filterValue)
			if !val.IsValid() || (val.Kind() == reflect.Ptr && val.IsNil()) || (val.Kind() == reflect.Interface && val.IsNil()) {
				continue
			}

			condition := fmt.Sprintf("%s = ?", field)
			db = db.Where(condition, filterValue)
		}

		return db
	}
}

func PreloadScope(relations ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, rel := range relations {
			db = db.Preload(rel)
		}
		return db
	}
}
