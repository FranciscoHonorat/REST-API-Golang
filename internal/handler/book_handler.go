package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/FranciscoHonorat/books-api/internal/domain"
	"github.com/FranciscoHonorat/books-api/internal/repository"
	"github.com/FranciscoHonorat/books-api/internal/service"
	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	service service.BookServicer
}

type CreateBookRequest struct {
	Name        string `json:"name" binding:"required"`
	Author      string `json:"author" binding:"required"`
	Year        int    `json:"year" binding:"required"`
	Masterpiece bool   `json:"masterpiece"`
}

type ListBooksResponse struct {
	Data  []domain.Book `json:"data"`
	Page  int32         `json:"page"`
	Limit int32         `json:"limit"`
	Total int32         `json:"total"`
}

func NewBookHandler(service service.BookServicer) *BookHandler {
	return &BookHandler{service: service}
}

func (h *BookHandler) GetBookByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	book, err := h.service.GetBook(c.Request.Context(), idInt)

	if errors.Is(err, domain.ErrBookNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) ListBooks(c *gin.Context) {
	author := c.Query("author")
	year := c.Query("year")
	masterpiece := c.Query("masterpiece")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	sort := c.DefaultQuery("sort", "name")

	authorStr := author

	var yearInt32 *int32
	if year != "" {
		yearInt, err := strconv.Atoi(year)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
			return
		}
		val := int32(yearInt)
		yearInt32 = &val
	}

	var masterpiecePtr *bool
	if masterpiece != "" {
		val, err := strconv.ParseBool(masterpiece)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid masterpiece value"})
			return
		}
		masterpiecePtr = &val
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page (must be >= 1)"})
		return
	}
	pageInt32 := int32(pageInt)

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 || limitInt > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit (must be between 1 and 100)"})
		return
	}

	validSorts := map[string]bool{"name": true, "author": true, "masterpiece": true}
	if !validSorts[sort] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort (must be name, author, or masterpiece)"})
		return
	}

	var authorPtr *string
	if author != "" {
		authorPtr = &authorStr
	}

	filter := repository.ListFilters{
		Author:      authorPtr,
		Year:        yearInt32,
		Masterpiece: masterpiecePtr,
	}

	pagination := repository.Pagination{
		Page:  pageInt32,
		Limit: int32(limitInt),
	}

	sorting := repository.Sorting{
		By: sort,
	}

	books, total, err := h.service.ListBooksWithTotal(c.Request.Context(), filter, pagination, sorting)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	response := ListBooksResponse{
		Data:  books,
		Page:  pageInt32,
		Limit: int32(limitInt),
		Total: total,
	}

	c.JSON(http.StatusOK, response)
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var request CreateBookRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON or missing required fields"})
		return
	}

	bookInput := &domain.Book{
		Name:        request.Name,
		Author:      request.Author,
		Year:        request.Year,
		Masterpiece: request.Masterpiece,
	}

	createdBook, err := h.service.CreateBook(c.Request.Context(), bookInput)

	if errors.Is(err, domain.ErrInvalidBookData) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book data"})
		return
	}

	if err != nil {
		slog.Error("CreateBook error", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, createdBook)
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var bookInput domain.Book
	if err := c.ShouldBindJSON(&bookInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	bookInput.ID = idInt

	updatedBook, err := h.service.UpdateBook(c.Request.Context(), &bookInput)

	if errors.Is(err, domain.ErrBookNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, updatedBook)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	idInt, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.service.DeleteBook(c.Request.Context(), idInt)

	if errors.Is(err, domain.ErrBookNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
