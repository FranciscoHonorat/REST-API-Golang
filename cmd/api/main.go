package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FranciscoHonorat/books-api/internal/config"
	"github.com/FranciscoHonorat/books-api/internal/handler"
	"github.com/FranciscoHonorat/books-api/internal/infra/postgres"
	"github.com/FranciscoHonorat/books-api/internal/middleware"
	"github.com/FranciscoHonorat/books-api/internal/seed"
	"github.com/FranciscoHonorat/books-api/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	setupLogger(cfg.Logging.Level)

	db, err := postgres.NewConnection(cfg.Database.URL)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	repo := postgres.NewPostgresBookRepository(db.Queries())
	svc := service.NewBookService(repo)
	bookHandler := handler.NewBookHandler(svc)

	if err := seed.SeedDatabase(svc, "books.json"); err != nil {
		log.Printf("Warning: seed failed: %v\n", err)
	}

	r := gin.New()

	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.RateLimitMiddleware())
	r.Use(middleware.MetricsMiddleware())
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.ValidationMiddleware())
	r.Use(middleware.CorsConfig())
	r.Use(middleware.LoggingMiddleware())

	v1 := r.Group("/api/v1")
	v1.GET("/books/:id", bookHandler.GetBookByID)
	v1.GET("/books", bookHandler.ListBooks)
	v1.POST("/books", bookHandler.CreateBook)
	v1.PUT("/books/:id", bookHandler.UpdateBook)
	v1.DELETE("/books/:id", bookHandler.DeleteBook)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)
		<-sigterm
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server forced to shutdown: %v\n", err)
		}
	}()

	log.Printf("Starting server on %s (TLS: %v)\n", server.Addr, cfg.Server.TLS)
	if cfg.Server.TLS {
		err = server.ListenAndServeTLS("certs/cert.pem", "certs/key.pem")
	} else {
		err = server.ListenAndServe()
	}
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func setupLogger(level string) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))
}
