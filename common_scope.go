package gormhelper

import (
	"fmt"
	"reflect"
	"strings"
	"the-core/core/authutil"
	"the-core/core/http/request"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func PaginationScope(req request.PaginationRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Limit(req.Size).
			Offset(req.Page - 1)
	}
}

func ActiveStatusScope(req request.IsAktifRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = ?", req.IsAktif)
	}
}

func IsDeletedScope(req request.IsDeleted) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if req.IsDeleted == nil {
			return db.Unscoped()
		}

		if *req.IsDeleted == true {
			return db.Unscoped().Where("deleted_at IS NOT NULL")
		}

		return db
	}
}

func SearchScope(req request.SearchRequest, fieldsToSearch ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if req.Search == "" || len(fieldsToSearch) == 0 {
			return db
		}

		pattern := "%" + req.Search + "%"

		var conditions []string
		var args []any

		for _, field := range fieldsToSearch {
			conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", field))
			args = append(args, pattern)
		}

		orConditions := strings.Join(conditions, " OR ")
		return db.Where(orConditions, args...)
	}
}

func SelfUserId(c fiber.Ctx) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", authutil.GetUserId(c))
	}
}

func FieldFilterScope(fieldFilters map[string]any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for field, filterValue := range fieldFilters {

			val := reflect.ValueOf(filterValue)
			if !val.IsValid() || (val.Kind() == reflect.Ptr && val.IsNil()) || (val.Kind() == reflect.Interface && val.IsNil()) {
				// fmt.Printf("field filter for %s is nil, skipping\n", field)
				continue
			}

			condition := fmt.Sprintf("%s = ?", field)
			db.Where(condition, filterValue)
		}

		return db
	}
}

func PreloadScope(modelsToPreload ...string) func(queryChain *gorm.DB) *gorm.DB {
	return func(queryChain *gorm.DB) *gorm.DB {
		for _, model := range modelsToPreload {
			queryChain.Preload(model)
		}

		return queryChain
	}
}
