package config

import (
	"os"
)

type Config struct {
	ServerAddress string
	ServerPort    string
	// TODO: Add more configuration fields as needed
}

func Load() (*Config, error) {
	return &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", "localhost"),
		ServerPort:    getEnv("SERVER_PORT", "8080"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
