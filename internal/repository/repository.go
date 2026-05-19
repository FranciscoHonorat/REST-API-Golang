package repository

import (
	"context"

	"github.com/FranciscoHonorat/books-api/internal/domain"
)

// BookRepository is the interface that defines the methods for interacting with the book repository.
type BookRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Book, error)
	List(ctx context.Context, filters ListFilters, pagination Pagination, sorting Sorting) ([]domain.Book, error)
	Count(ctx context.Context, filters ListFilters) (int32, error)
	Create(ctx context.Context, book *domain.Book) (*domain.Book, error)
	Update(ctx context.Context, book *domain.Book) (*domain.Book, error)
	Delete(ctx context.Context, id int64) error
}

// ListFilters represents the filters that can be applied when listing books.
type ListFilters struct {
	Author      *string
	Year        *int32
	Masterpiece *bool
}

// Pagination represents the pagination parameters for listing books.
type Pagination struct {
	Page  int32
	Limit int32
}

// Sorting represents the sorting parameters for listing books.
type Sorting struct {
	By string // "name", "author", "masterpiece"
}
