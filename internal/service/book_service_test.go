package service

import (
	"context"
	"errors"
	"testing"

	"github.com/FranciscoHonorat/books-api/internal/domain"
	"github.com/FranciscoHonorat/books-api/internal/repository"
	"github.com/stretchr/testify/assert"
)

// MockBookRepository for testing
type MockBookRepository struct {
	GetByIDFunc func(ctx context.Context, id int64) (*domain.Book, error)
	ListFunc    func(ctx context.Context, filters repository.ListFilters, pagination repository.Pagination, sorting repository.Sorting) ([]domain.Book, error)
	CountFunc   func(ctx context.Context, filters repository.ListFilters) (int32, error)
	CreateFunc  func(ctx context.Context, book *domain.Book) (*domain.Book, error)
	UpdateFunc  func(ctx context.Context, book *domain.Book) (*domain.Book, error)
	DeleteFunc  func(ctx context.Context, id int64) error
}

func (m *MockBookRepository) GetByID(ctx context.Context, id int64) (*domain.Book, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m *MockBookRepository) List(ctx context.Context, filters repository.ListFilters, pagination repository.Pagination, sorting repository.Sorting) ([]domain.Book, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, filters, pagination, sorting)
	}
	return []domain.Book{}, nil
}

func (m *MockBookRepository) Count(ctx context.Context, filters repository.ListFilters) (int32, error) {
	if m.CountFunc != nil {
		return m.CountFunc(ctx, filters)
	}
	return 0, nil
}

func (m *MockBookRepository) Create(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, book)
	}
	return book, nil
}

func (m *MockBookRepository) Update(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, book)
	}
	return book, nil
}

func (m *MockBookRepository) Delete(ctx context.Context, id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func TestCreateBook_Success(t *testing.T) {
	// ARRANGE
	mockRepo := &MockBookRepository{
		CreateFunc: func(ctx context.Context, book *domain.Book) (*domain.Book, error) {
			book.ID = 1 // Simulate ID assignment
			return book, nil
		},
	}

	svc := NewBookService(mockRepo)
	book := &domain.Book{
		Name:   "Test Book",
		Author: "Test Author",
		Year:   2020,
	}

	// ACT
	createdBook, err := svc.CreateBook(context.Background(), book)

	// ASSERT
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if createdBook.ID != 1 {
		t.Errorf("Expected book ID to be 1, got %d", createdBook.ID)
	}
	if createdBook.Name != book.Name {
		t.Errorf("Expected book name to be '%s', got '%s'", book.Name, createdBook.Name)
	}
}

func TestCreateBook_InvalidData(t *testing.T) {
	// ARRANGE
	mockRepo := &MockBookRepository{
		CreateFunc: func(ctx context.Context, book *domain.Book) (*domain.Book, error) {
			return nil, domain.ErrInvalidBookData
		},
	}

	svc := NewBookService(mockRepo)
	book := &domain.Book{
		Name:   "", // Invalid name
		Author: "Test Author",
		Year:   2020,
	}

	// ACT
	_, err := svc.CreateBook(context.Background(), book)

	// ASSERT
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	assert.True(t, errors.Is(err, domain.ErrInvalidBookData))

	if err.Error() != "invalid book data" {
		t.Errorf("Expected error message to be 'invalid book data', got '%s'", err.Error())
	}
}

func TestGetBook_Success(t *testing.T) {
	// ARRANGE
	mockBook := &domain.Book{
		ID:     1,
		Name:   "Test Book",
		Author: "Test Author",
		Year:   2020,
	}

	mockRepo := &MockBookRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*domain.Book, error) {
			return mockBook, nil
		},
	}

	svc := NewBookService(mockRepo)

	// ACT
	foundBook, err := svc.GetBook(context.Background(), 1)

	// ASSERT
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if foundBook.ID != mockBook.ID {
		t.Errorf("Expected book ID to be %d, got %d", mockBook.ID, foundBook.ID)
	}
}

func TestGetBook_NotFound(t *testing.T) {
	// ARRANGE
	mockRepo := &MockBookRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*domain.Book, error) {
			return nil, domain.ErrBookNotFound
		},
	}

	svc := NewBookService(mockRepo)

	// ACT
	_, err := svc.GetBook(context.Background(), 999)

	// ASSERT
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	assert.True(t, errors.Is(err, domain.ErrBookNotFound))
}

func TestListBooks_Success(t *testing.T) {
	// ARRANGE
	mockBooks := []domain.Book{
		{ID: 1, Name: "Book 1", Author: "Author 1", Year: 2020},
		{ID: 2, Name: "Book 2", Author: "Author 2", Year: 2021},
	}

	mockRepo := &MockBookRepository{
		ListFunc: func(ctx context.Context, filters repository.ListFilters, pagination repository.Pagination, sorting repository.Sorting) ([]domain.Book, error) {
			return mockBooks, nil
		},
	}
	svc := NewBookService(mockRepo)

	// ACT
	books, err := svc.ListBooks(context.Background(), repository.ListFilters{}, repository.Pagination{}, repository.Sorting{})

	// ASSERT
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(books) != len(mockBooks) {
		t.Errorf("Expected %d books, got %d", len(mockBooks), len(books))
	}
}
