package domain

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Pagination struct {
	Page     int `validate:"gte=0"`
	PageSize int `validate:"gte=0,lte=100"`
	Total    int
}

func NewPagination(page, pageSize, total int) *Pagination {
	return &Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}
}

func (p *Pagination) Validate() error {
	pagValidator := validator.New()

	if err := pagValidator.Struct(p); err != nil {
		return fmt.Errorf("paginations validate struct: %w", err)
	}
	return nil
}
