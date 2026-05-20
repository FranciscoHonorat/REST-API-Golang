package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/FranciscoHonorat/books-api/internal/handler"
	"github.com/FranciscoHonorat/books-api/internal/infra/postgres"
	"github.com/FranciscoHonorat/books-api/internal/middleware"
	"github.com/FranciscoHonorat/books-api/internal/seed"
	"github.com/FranciscoHonorat/books-api/internal/service"
)

func setupE2EEngine() *gin.Engine {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}
	if dsn == "" {
		panic("TEST_DATABASE_URL environment variable is not set")
	}

	db, err := postgres.NewConnection(dsn)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect: %v", err))
	}

	queries := db.Queries()
	bookRepo := postgres.NewPostgresBookRepository(queries)

	bookService := service.NewBookService(bookRepo)

	err = seed.SeedDatabase(bookService, "../books.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to seed: %v", err))
	}

	bookHandler := handler.NewBookHandler(bookService)

	gin.SetMode(gin.TestMode)
	engine := gin.New()

	engine.Use(middleware.RecoveryMiddleware())

	engine.POST("/books", bookHandler.CreateBook)
	engine.GET("/books", bookHandler.ListBooks)
	engine.GET("/books/:id", bookHandler.GetBookByID)
	engine.PUT("/books/:id", bookHandler.UpdateBook)
	engine.DELETE("/books/:id", bookHandler.DeleteBook)

	return engine
}

func TestCreateBook_E2E(t *testing.T) {
	engine := setupE2EEngine()

	payload := map[string]interface{}{
		"name":        "Cien años de soledad",
		"author":      "Gabriel García Márquez",
		"year":        1967,
		"masterpiece": true,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}
	req, err := http.NewRequest("POST", "/books", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["id"] == nil {
		t.Fatalf("Expected 'id' field in response, got nil")
	}
	assert.Equal(t, "Cien años de soledad", response["name"])
	assert.NotZero(t, response["id"])
}

func TestListBooksE2E(t *testing.T) {
	engine := setupE2EEngine()
	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if len(response) > 0 {
		firstBook := response[0]

		assert.NotNil(t, firstBook["id"])
		assert.NotNil(t, firstBook["name"])
		assert.NotNil(t, firstBook["author"])
		assert.NotNil(t, firstBook["year"])
		assert.NotNil(t, firstBook["masterpiece"])
		assert.NotNil(t, firstBook["created_at"])
		assert.NotNil(t, firstBook["updated_at"])
	}
}

func TestGetBookByID_E2E(t *testing.T) {
	engine := setupE2EEngine()
	req, err := http.NewRequest("GET", "/books/1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["id"] == nil {
		t.Fatalf("Expected 'id' field in response, got nil")
	}
	assert.Equal(t, float64(1), response["id"])
	assert.NotNil(t, response["name"])
	assert.NotNil(t, response["author"])
	assert.NotNil(t, response["year"])
	assert.NotNil(t, response["masterpiece"])
	assert.NotNil(t, response["created_at"])
	assert.NotNil(t, response["updated_at"])

	_, nameOk := response["name"].(string)
	assert.True(t, nameOk, "Expected 'name' to be a string")

	_, authorOk := response["author"].(string)
	assert.True(t, authorOk, "Expected 'author' to be a string")
}

func TestUpdateBook_E2E(t *testing.T) {
	engine := setupE2EEngine()
	payload := map[string]interface{}{
		"name":        "Cien años de soledad (Actualizado)",
		"author":      "Gabriel García Márquez",
		"year":        1967,
		"masterpiece": true,
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("PUT", "/books/1", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["id"] == nil {
		t.Fatalf("Expected 'id' field in response, got nil")
	}
	assert.Equal(t, "Cien años de soledad (Actualizado)", response["name"])
}

func TestDeleteBook_E2E(t *testing.T) {
	engine := setupE2EEngine()
	req, err := http.NewRequest("DELETE", "/books/1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	req, err = http.NewRequest("GET", "/books/1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
