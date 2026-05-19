package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FranciscoHonorat/books-api/internal/domain"
	"github.com/FranciscoHonorat/books-api/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockBookRepository for testing
type MockBookService struct {
	GetBookFunc            func(ctx context.Context, id int64) (*domain.Book, error)
	ListBooksFunc          func(ctx context.Context, filters repository.ListFilters, pagination repository.Pagination, sorting repository.Sorting) ([]domain.Book, error)
	ListBooksWithTotalFunc func(ctx context.Context, filters repository.ListFilters, pagination repository.Pagination, sorting repository.Sorting) ([]domain.Book, int32, error)
	CreateBookFunc         func(ctx context.Context, book *domain.Book) (*domain.Book, error)
	UpdateBookFunc         func(ctx context.Context, book *domain.Book) (*domain.Book, error)
	DeleteBookFunc         func(ctx context.Context, id int64) error
}

func (m *MockBookService) GetBook(ctx context.Context, id int64) (*domain.Book, error) {
	return m.GetBookFunc(ctx, id)
}

func (m *MockBookService) ListBooks(ctx context.Context, filters repository.ListFilters, pagination repository.Pagination, sorting repository.Sorting) ([]domain.Book, error) {
	return []domain.Book{}, nil
}

func (m *MockBookService) ListBooksWithTotal(ctx context.Context, filters repository.ListFilters, pagination repository.Pagination, sorting repository.Sorting) ([]domain.Book, int32, error) {
	if m.ListBooksWithTotalFunc != nil {
		return m.ListBooksWithTotalFunc(ctx, filters, pagination, sorting)
	}
	return []domain.Book{}, 0, nil
}

func (m *MockBookService) CreateBook(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	if m.CreateBookFunc != nil {
		return m.CreateBookFunc(ctx, book)
	}
	return book, nil
}

func (m *MockBookService) UpdateBook(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	if m.UpdateBookFunc != nil {
		return m.UpdateBookFunc(ctx, book)
	}
	return book, nil
}

func (m *MockBookService) DeleteBook(ctx context.Context, id int64) error {
	if m.DeleteBookFunc != nil {
		return m.DeleteBookFunc(ctx, id)
	}
	return nil
}

func TestGetBookByID(t *testing.T) {
	// ARRANGE
	mockBook := &domain.Book{
		ID:     1,
		Name:   "Test Book",
		Author: "Test Author",
		Year:   2020,
	}

	mockService := &MockBookService{
		GetBookFunc: func(ctx context.Context, id int64) (*domain.Book, error) {
			return mockBook, nil
		},
	}

	handler := NewBookHandler(mockService)

	// ACT
	req := httptest.NewRequest("GET", "/books/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	handler.GetBookByID(c)

	// ASSERT
	assert.Equal(t, http.StatusOK, w.Code)

	var responseBook domain.Book
	err := json.Unmarshal(w.Body.Bytes(), &responseBook)
	assert.NoError(t, err)
	assert.Equal(t, mockBook.ID, responseBook.ID)
	assert.Equal(t, mockBook.Name, responseBook.Name)
}

func TestGetBookByID_NotFound(t *testing.T) {
	// ARRANGE
	mockService := &MockBookService{
		GetBookFunc: func(ctx context.Context, id int64) (*domain.Book, error) {
			return nil, domain.ErrBookNotFound
		},
	}

	handler := NewBookHandler(mockService)

	// ACT
	req := httptest.NewRequest("GET", "/books/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	handler.GetBookByID(c)

	// ASSERT
	assert.Equal(t, http.StatusNotFound, w.Code)

	var errorResponse map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Book not found", errorResponse["error"])
}

// TestCreateBook
func TestCreateBook(t *testing.T) {
	// ARRANGE
	mockService := &MockBookService{
		CreateBookFunc: func(ctx context.Context, book *domain.Book) (*domain.Book, error) {
			book.ID = 1 // Simulate ID assignment
			return book, nil
		},
	}

	handler := NewBookHandler(mockService)

	// ACT
	bookJSON := `{"name":"New Book","author":"New Author","year":2021}`
	req := httptest.NewRequest("POST", "/books", nil)
	req.Body = io.NopCloser(strings.NewReader(bookJSON))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler.CreateBook(c)

	// ASSERT
	assert.Equal(t, http.StatusCreated, w.Code)

	var responseBook domain.Book
	err := json.Unmarshal(w.Body.Bytes(), &responseBook)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), responseBook.ID)
	assert.Equal(t, "New Book", responseBook.Name)
	assert.Equal(t, "New Author", responseBook.Author)
	assert.Equal(t, 2021, responseBook.Year)
}

// TestUpdateBook
func TestUpdateBook(t *testing.T) {
	// ARRANGE
	mockService := &MockBookService{
		UpdateBookFunc: func(ctx context.Context, book *domain.Book) (*domain.Book, error) {
			return book, nil
		},
	}

	handler := NewBookHandler(mockService)

	// ACT
	bookJSON := `{"name":"Updated Book","author":"Updated Author","year":2022}`
	req := httptest.NewRequest("PUT", "/books/1", nil)
	req.Body = io.NopCloser(strings.NewReader(bookJSON))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	handler.UpdateBook(c)

	// ASSERT
	assert.Equal(t, http.StatusOK, w.Code)

	var responseBook domain.Book
	err := json.Unmarshal(w.Body.Bytes(), &responseBook)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), responseBook.ID)
	assert.Equal(t, "Updated Book", responseBook.Name)
	assert.Equal(t, "Updated Author", responseBook.Author)
	assert.Equal(t, 2022, responseBook.Year)
}

// TestDeleteBook
func TestDeleteBook(t *testing.T) {
	// ARRANGE
	mockService := &MockBookService{
		DeleteBookFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	handler := NewBookHandler(mockService)
	router := gin.New()
	router.DELETE("/books/:id", handler.DeleteBook)

	// ACT
	req := httptest.NewRequest("DELETE", "/books/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ASSERT
	assert.Equal(t, http.StatusNoContent, w.Code)
}
