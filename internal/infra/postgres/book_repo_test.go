package postgres

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/FranciscoHonorat/books-api/internal/domain"
	"github.com/FranciscoHonorat/books-api/internal/infra/sqlc"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a test database connection
func setupTestDB(t *testing.T) *sql.DB {
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		testDBURL = "postgres://postgres:postgres@localhost:5432/books_test_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", testDBURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	db.ExecContext(context.Background(), "TRUNCATE TABLE books")
	return db
}

func TestCreateBook_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupTestDB(t)
	if db == nil {
		t.Skip("Test database not available")
	}
	defer db.Close()

	queries := sqlc.New(db)
	repo := NewPostgresBookRepository(queries)

	book := &domain.Book{
		Name:        "Integration Test Book",
		Author:      "Integration Test Author",
		Year:        2021,
		Masterpiece: true,
	}

	createdBook, err := repo.Create(context.Background(), book)

	assert.NoError(t, err)
	assert.NotNil(t, createdBook)
	assert.Equal(t, book.Name, createdBook.Name)
	assert.Equal(t, book.Author, createdBook.Author)
}

func TestGetByID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupTestDB(t)
	if db == nil {
		t.Skip("Test database not available")
	}
	defer db.Close()

	queries := sqlc.New(db)
	repo := NewPostgresBookRepository(queries)

	book := &domain.Book{
		Name:        "Integration Test Book",
		Author:      "Integration Test Author",
		Year:        2021,
		Masterpiece: false,
	}
	created, _ := repo.Create(context.Background(), book)

	found, err := repo.GetByID(context.Background(), created.ID)

	assert.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
	assert.Equal(t, created.Name, found.Name)
}

func TestDelete_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupTestDB(t)
	if db == nil {
		t.Skip("Test database not available")
	}
	defer db.Close()

	queries := sqlc.New(db)
	repo := NewPostgresBookRepository(queries)

	book := &domain.Book{
		Name:        "Integration Test Book",
		Author:      "Integration Test Author",
		Year:        2021,
		Masterpiece: false,
	}
	created, _ := repo.Create(context.Background(), book)

	err := repo.Delete(context.Background(), created.ID)

	assert.NoError(t, err)

	_, err = repo.GetByID(context.Background(), created.ID)
	assert.True(t, errors.Is(err, domain.ErrBookNotFound))
}
