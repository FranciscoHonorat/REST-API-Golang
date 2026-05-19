package domain

import (
	"strings"
	"time"
)

type Book struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Author      string    `json:"author"`
	Year        int       `json:"year"`
	Masterpiece bool      `json:"masterpiece"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (b *Book) Validate() error {
	if strings.TrimSpace(b.Name) == "" {
		return ErrInvalidBookData
	}

	if strings.TrimSpace(b.Author) == "" {
		return ErrInvalidBookData
	}

	if b.Year <= 0 || b.Year > time.Now().Year() {
		return ErrInvalidBookData
	}

	return nil
}

func NewBook(name, author string, year int, masterpiece bool) (*Book, error) {
	book := &Book{
		Name:        name,
		Author:      author,
		Year:        year,
		Masterpiece: masterpiece,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := book.Validate(); err != nil {
		return nil, err
	}

	return book, nil
}
