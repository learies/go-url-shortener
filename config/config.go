package config

import (
	"flag"
)

const (
	DefaultAddress = ":8080"
	DefaultBaseURL = "http://localhost:8080"
)

// Config структура, которая хранит аргументы командной строки
type Config struct {
	Address string
	BaseURL string
}

// ParseConfig функция для разбора аргументов командной строки и возврата настройки
func ParseConfig() *Config {
	addressFlag := flag.String("a", DefaultAddress, "Address to start the HTTP server")
	baseURLFlag := flag.String("b", DefaultBaseURL, "Base URL for the shortened URL")

	flag.Parse()

	return &Config{
		Address: *addressFlag,
		BaseURL: *baseURLFlag,
	}
}
