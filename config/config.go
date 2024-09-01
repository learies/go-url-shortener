package config

import (
	"flag"
)

// Config структура, которая хранит аргументы командной строки
type Config struct {
	Address string
	BaseURL string
}

// ParseConfig функция для разбора аргументов командной строки и возврата настройки
func ParseConfig() *Config {
	address := flag.String("a", "localhost:8080", "Address to start the HTTP server")
	baseURL := flag.String("b", "http://localhost:8080", "Base URL for the shortened URL")

	flag.Parse()

	return &Config{
		Address: *address,
		BaseURL: *baseURL,
	}
}
