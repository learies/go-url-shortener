package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

const (
	DefaultAddress = ":8080"
	DefaultBaseURL = "http://localhost:8080"
)

// Config структура, которая хранит аргументы командной строки и переменные окружения
type Config struct {
	Address string `env:"SERVER_ADDRESS"`
	BaseURL string `env:"BASE_URL"`
}

// ParseConfig функция для разбора аргументов командной строки и переменных окружения
func ParseConfig() *Config {
	cfg := &Config{
		Address: DefaultAddress,
		BaseURL: DefaultBaseURL,
	}

	// Парсинг переменных окружения
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}

	// Парсинг флагов командной строки
	addressFlag := flag.String("a", cfg.Address, "Address to start the HTTP server")
	baseURLFlag := flag.String("b", cfg.BaseURL, "Base URL for the shortened URL")
	flag.Parse()

	// Приоритизация: если флаги командной строки не заданы, оставляем значения из env
	if addressFlag != nil && *addressFlag != DefaultAddress {
		cfg.Address = *addressFlag
	}

	if baseURLFlag != nil && *baseURLFlag != DefaultBaseURL {
		cfg.BaseURL = *baseURLFlag
	}

	return cfg
}
