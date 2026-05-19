package seed

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/FranciscoHonorat/books-api/internal/domain"
	"github.com/FranciscoHonorat/books-api/internal/repository"
	"github.com/FranciscoHonorat/books-api/internal/service"
)

func SeedDatabase(svc service.BookServicer, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	var books []domain.Book
	if err := json.Unmarshal(data, &books); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}

	existingBooks, err := svc.ListBooks(context.Background(),
		repository.ListFilters{},
		repository.Pagination{Page: 1, Limit: 1},
		repository.Sorting{By: "name"})

	if err == nil && len(existingBooks) > 0 {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	for _, book := range books {
		if _, err := svc.CreateBook(context.Background(), &book); err != nil {
			log.Printf("Warning: failed to create book '%s': %v\n", book.Name, err)
			continue
		}
	}

	log.Println("Database seeded successfully!")
	return nil
}
