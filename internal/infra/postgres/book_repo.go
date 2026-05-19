package postgres

import (
	"context"
	"database/sql"

	"github.com/FranciscoHonorat/books-api/internal/domain"
	"github.com/FranciscoHonorat/books-api/internal/infra/sqlc"
	"github.com/FranciscoHonorat/books-api/internal/repository"
)

// PostgresBookRepository implements the BookRepository interface using PostgreSQL
type PostgresBookRepository struct {
	queries *sqlc.Queries
}

// NewPostgresBookRepository creates a new instance of the repository.
func NewPostgresBookRepository(queries *sqlc.Queries) *PostgresBookRepository {
	return &PostgresBookRepository{
		queries: queries,
	}
}

// GetByID searches for a book by ID.
func (r *PostgresBookRepository) GetByID(ctx context.Context, id int64) (*domain.Book, error) {
	sqlcBook, err := r.queries.GetBookByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrBookNotFound
		}
		return nil, domain.ErrInternalServer
	}
	return r.convertSqlcBookToDomain(&sqlcBook), nil
}

// List lists books with filters, pagination, and sorting.
func (r *PostgresBookRepository) List(ctx context.Context, filters repository.ListFilters, pagination repository.Pagination, sorting repository.Sorting) ([]domain.Book, error) {
	params := sqlc.ListBooksParams{
		PageOffset: (pagination.Page - 1) * pagination.Limit,
		PageLimit:  pagination.Limit,
		SortBy:     sorting.By,
	}

	if filters.Author != nil {
		params.FilterAuthor = sql.NullString{String: *filters.Author, Valid: true}
	}

	if filters.Year != nil {
		params.FilterYear = sql.NullInt32{Int32: *filters.Year, Valid: true}
	}

	if filters.Masterpiece != nil {
		params.FilterMasterpiece = sql.NullBool{Bool: *filters.Masterpiece, Valid: true}
	}

	sqlcBooks, err := r.queries.ListBooks(ctx, params)
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	domainBooks := make([]domain.Book, len(sqlcBooks))
	for i, sqlcBook := range sqlcBooks {
		domainBooks[i] = *r.convertSqlcBookToDomain(&sqlcBook)
	}
	return domainBooks, nil
}

func (r *PostgresBookRepository) Count(ctx context.Context, filters repository.ListFilters) (int32, error) {
	params := sqlc.CountBooksParams{
		FilterAuthor:      sql.NullString{},
		FilterYear:        sql.NullInt32{},
		FilterMasterpiece: sql.NullBool{},
	}

	if filters.Author != nil {
		params.FilterAuthor = sql.NullString{String: *filters.Author, Valid: true}
	}

	if filters.Year != nil {
		params.FilterYear = sql.NullInt32{Int32: *filters.Year, Valid: true}
	}

	if filters.Masterpiece != nil {
		params.FilterMasterpiece = sql.NullBool{Bool: *filters.Masterpiece, Valid: true}
	}

	count, err := r.queries.CountBooks(ctx, params)
	if err != nil {
		return 0, domain.ErrInternalServer
	}

	return int32(count), nil
}

func (r *PostgresBookRepository) Create(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	params := sqlc.CreateBookParams{
		Name:        book.Name,
		Author:      book.Author,
		Year:        int32(book.Year),
		Masterpiece: book.Masterpiece,
	}

	sqlcBook, err := r.queries.CreateBook(ctx, params)
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	return r.convertSqlcBookToDomain(&sqlcBook), nil
}

func (r *PostgresBookRepository) Update(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	params := sqlc.UpdateBookParams{
		ID:          book.ID,
		Name:        book.Name,
		Author:      book.Author,
		Year:        int32(book.Year),
		Masterpiece: book.Masterpiece,
	}

	sqlcBook, err := r.queries.UpdateBook(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrBookNotFound
		}
		return nil, domain.ErrInternalServer
	}

	return r.convertSqlcBookToDomain(&sqlcBook), nil
}

func (r *PostgresBookRepository) Delete(ctx context.Context, id int64) error {
	rowsAffected, err := r.queries.DeleteBook(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrBookNotFound
	}
	return nil
}

func (r *PostgresBookRepository) convertSqlcBookToDomain(sqlcBook *sqlc.Book) *domain.Book {
	return &domain.Book{
		ID:          sqlcBook.ID,
		Name:        sqlcBook.Name,
		Author:      sqlcBook.Author,
		Year:        int(sqlcBook.Year),
		Masterpiece: sqlcBook.Masterpiece,
		CreatedAt:   sqlcBook.CreatedAt,
		UpdatedAt:   sqlcBook.UpdatedAt,
	}
}
