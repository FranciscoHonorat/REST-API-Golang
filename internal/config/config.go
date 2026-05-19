package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	maxConnsStr := os.Getenv("DATABASE_MAX_CONNECTIONS")
	maxConns := 25
	if maxConnsStr != "" {
		maxConnsParsed, err := strconv.Atoi(maxConnsStr)
		if err != nil || maxConnsParsed <= 0 {
			return nil, errors.New("DATABASE_MAX_CONNECTIONS must be a positive integer")
		}
		maxConns = maxConnsParsed
	}

	sslMode := os.Getenv("DATABASE_SSL_MODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	portStr := os.Getenv("SERVER_PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 || port > 65535 {
		return nil, errors.New("SERVER_PORT must be a valid port number (1-65535)")
	}

	env := os.Getenv("SERVER_ENV")
	if env != "development" && env != "production" && env != "staging" {
		return nil, errors.New("SERVER_ENV must be one of 'development', 'production', or 'staging'")
	}

	tlsStr := os.Getenv("SERVER_TLS")
	tls := env == "production"
	if tlsStr != "" {
		tlsParsed, err := strconv.ParseBool(tlsStr)
		if err == nil {
			tls = tlsParsed
		}
	}

	level := os.Getenv("LOGGING_LEVEL")
	if level == "" {
		level = "info"
	}
	if level != "debug" && level != "info" && level != "warn" && level != "error" {
		return nil, errors.New("LOGGING_LEVEL must be one of 'debug', 'info', 'warn', or 'error'")
	}

	format := os.Getenv("LOGGING_FORMAT")
	if format == "" {
		format = "json"
	}

	return &Config{
		Database: DatabaseConfig{
			URL:            dbURL,
			MaxConnections: maxConns,
			SSLMode:        sslMode,
		},
		Server: ServerConfig{
			Port: port,
			Env:  env,
			TLS:  tls,
		},
		Logging: LoggingConfig{
			Level:  level,
			Format: format,
		},
	}, nil
}

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Logging  LoggingConfig
}

type DatabaseConfig struct {
	URL            string
	MaxConnections int
	SSLMode        string
}

type ServerConfig struct {
	Port int
	Env  string
	TLS  bool
}

type LoggingConfig struct {
	Level  string
	Format string
}
