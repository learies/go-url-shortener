package config

import (
	"flag"
	"os"
)

type Config struct {
	Address string
	BaseURL string
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func LoadConfig() Config {
	// Значения по умолчанию
	defaultAddress := "localhost:8080"
	defaultBaseURL := "http://localhost:8080"

	// Чтение из переменных окружения
	envAddress := getEnv("SERVER_ADDRESS", defaultAddress)
	envBaseURL := getEnv("BASE_URL", defaultBaseURL)

	// Чтение из флагов командной строки
	address := flag.String("a", envAddress, "address to start the HTTP server")
	baseURL := flag.String("b", envBaseURL, "base URL for shortened URLs")
	flag.Parse()

	return Config{
		Address: *address,
		BaseURL: *baseURL,
	}
}
