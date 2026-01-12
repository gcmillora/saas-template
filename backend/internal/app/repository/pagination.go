package repository

import pg "github.com/go-jet/jet/v2/postgres"


type PaginationParams struct {
	Page int
	PageSize int
	SortBy string
	SortDesc bool
}

type PaginatedResponse[T any] struct {
	Data []T `json:"data"`
	TotalCount int `json:"totalCount"`
	Page int `json:"page"`
	PageSize int `json:"pageSize"`
	TotalPages int `json:"totalPages"`
}

func DefaultPaginationParams() PaginationParams {
	return PaginationParams{
		Page: 1,
		PageSize: 20,
	}
}

func (p PaginationParams) GetOffset() int64 {
	if p.Page < 1 {
		p.Page = 1
	}

	return int64((p.Page - 1) * p.PageSize)
}

func (p PaginationParams) GetLimit() int64 {
	if p.PageSize < 1 {
		p.PageSize = 20
	}

	return int64(p.PageSize)
}

func (p PaginationParams) ApplyPagination(stmt pg.SelectStatement) pg.SelectStatement {
	return stmt.LIMIT(p.GetLimit()).OFFSET(p.GetOffset())
}


func (p PaginationParams) CalculateTotalPages(totalCount int) int {
	if p.PageSize < 1 {
		p.PageSize = 20
	}
	totalPages := totalCount/ p.PageSize
	if totalCount % p.PageSize > 0 {
		totalPages++
	}

	return totalPages
}

func NewPaginatedResponse[T any](data []T, totalCount int, params PaginationParams) *PaginatedResponse[T] {
	return &PaginatedResponse[T]{
		Data: data,
		TotalCount: totalCount,
		Page: params.Page,
		PageSize: params.PageSize,
		TotalPages: params.CalculateTotalPages(totalCount),
	}
}