package service

import (
	"context"

	"github.com/FranciscoHonorat/books-api/internal/domain"
	"github.com/FranciscoHonorat/books-api/internal/repository"
)

type BookServicer interface {
	GetBook(ctx context.Context, id int64) (*domain.Book, error)
	ListBooks(ctx context.Context, filters repository.ListFilters, pagination repository.Pagination, sorting repository.Sorting) ([]domain.Book, error)
	ListBooksWithTotal(ctx context.Context, filters repository.ListFilters, pagination repository.Pagination, sorting repository.Sorting) ([]domain.Book, int32, error)
	CreateBook(ctx context.Context, book *domain.Book) (*domain.Book, error)
	UpdateBook(ctx context.Context, book *domain.Book) (*domain.Book, error)
	DeleteBook(ctx context.Context, id int64) error
}

type BookService struct {
	repo repository.BookRepository
}

// NewBookService()
func NewBookService(repo repository.BookRepository) *BookService {
	return &BookService{
		repo: repo,
	}
}

// GetBook()
func (s *BookService) GetBook(ctx context.Context, id int64) (*domain.Book, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidBookData
	}

	// Chamar com nome correto
	book, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return book, nil
}

// CreateBook()
func (s *BookService) CreateBook(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	// Validação completa
	if book == nil {
		return nil, domain.ErrInvalidBookData
	}
	if err := book.Validate(); err != nil {
		return nil, domain.ErrInvalidBookData
	}
	createdBook, err := s.repo.Create(ctx, book)
	if err != nil {
		return nil, err
	}
	return createdBook, nil
}

// ListBooks()
func (s *BookService) ListBooks(ctx context.Context, filters repository.ListFilters, pagination repository.Pagination, sorting repository.Sorting) ([]domain.Book, error) {
	//call repository
	books, err := s.repo.List(ctx, filters, pagination, sorting)
	if err != nil {
		return nil, err
	}

	//return result
	return books, nil
}

// ListBooksWithTotal()
func (s *BookService) ListBooksWithTotal(ctx context.Context, filters repository.ListFilters, pagination repository.Pagination, sorting repository.Sorting) ([]domain.Book, int32, error) {
	books, err := s.repo.List(ctx, filters, pagination, sorting)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

// UpdateBook()
func (s *BookService) UpdateBook(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	if book == nil {
		return nil, domain.ErrInvalidBookData
	}
	if err := book.Validate(); err != nil {
		return nil, domain.ErrInvalidBookData
	}
	updatedBook, err := s.repo.Update(ctx, book)
	if err != nil {
		return nil, err
	}

	//return result
	return updatedBook, nil
}

// DeleteBook()
func (s *BookService) DeleteBook(ctx context.Context, id int64) error {
	//validate input
	if id <= 0 {
		return domain.ErrBookNotFound
	}
	//calls repository
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
