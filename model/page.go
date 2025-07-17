package model

type Page[T any] struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	Items       []T `json:"items"`
	Total       int `json:"total"`
}

func NewPage[T any](currentPage, pageSize, total int, items []T) *Page[T] {
	return &Page[T]{
		CurrentPage: currentPage,
		PageSize:    pageSize,
		Items:       items,
		Total:       total,
	}
}
