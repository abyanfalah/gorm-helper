package github.com/abyanfalah/gorm-helper

import "github.com/google/uuid"

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
	IsDeleted *bool `json:"IsDeleted"`
}

type IsAktifRequest struct {
	IsAktif *bool `json:"isAktif"`
}

type SearchPaginationRequest struct {
	PaginationRequest
	SearchRequest
	IsDeleted
}
