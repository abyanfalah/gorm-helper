package gormhelper

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type IDNumberRequest struct {
	ID int `json:"id" validate:"required"`
}

type UUIDRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

type PaginationRequest struct {
	Page int `json:"page" validate:"required,min=1"`
	Size int `json:"size" validate:"required,min=1"`
}

type SearchRequest struct {
	Search string `json:"search"`
}

type IsDeleted struct {
	IsDeleted *bool `json:"isDeleted"`
}

type IsAktifRequest struct {
	IsAktif *bool `json:"isAktif"`
}

type SearchPaginationRequest struct {
	PaginationRequest
	SearchRequest
	IsDeleted
}

func (s *SearchPaginationRequest) SetFilters(c fiber.Ctx) {
	s.Page = fiber.Query(c, "page", 1)
	s.Size = fiber.Query(c, "size", 10)
	s.Search = strings.TrimSpace(c.Query("search"))

	if s.Page < 1 {
		s.Page = 1
	}
	if s.Size < 1 {
		s.Size = 10
	}
}