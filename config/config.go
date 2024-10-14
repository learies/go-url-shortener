package config

import (
	"flag"
	"os"
)

type Config struct {
	Address         string
	BaseURL         string
	FileStoragePath string
	DatabaseDSN     string
	LogLevel        string
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func LoadConfig() Config {
	// Default values
	defaultAddress := "localhost:8080"
	defaultBaseURL := "http://localhost:8080"
	var defaultFileStoragePath string
	defaultDatabaseDSN := "postgres://postgres:postgres@localhost/postgres?sslmode=disable"
	defaultLevel := "info"

	// Read from environment variables
	envAddress := getEnv("SERVER_ADDRESS", defaultAddress)
	envBaseURL := getEnv("BASE_URL", defaultBaseURL)
	envFileStoragePath := getEnv("FILE_STORAGE_PATH", defaultFileStoragePath)
	envDatabaseDSN := getEnv("DATABASE_DSN", defaultDatabaseDSN)
	envLevel := getEnv("LOG_LEVEL", defaultLevel)

	// Read from command-line flags
	address := flag.String("a", envAddress, "address to start the HTTP server")
	baseURL := flag.String("b", envBaseURL, "base URL for shortened URLs")
	fileStoragePath := flag.String("f", envFileStoragePath, "path to the file for storing URL data")
	databaseDSN := flag.String("d", envDatabaseDSN, "database DSN")
	logLevel := flag.String("l", envLevel, "log level")

	flag.Parse()

	return Config{
		Address:         *address,
		BaseURL:         *baseURL,
		FileStoragePath: *fileStoragePath,
		DatabaseDSN:     *databaseDSN,
		LogLevel:        *logLevel,
	}
}
