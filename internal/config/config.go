package config

import (
	"github.com/joho/godotenv"
	"github.com/sborsh1kmusora/auth/internal/logger"
)

func Load(path string) {
	if err := godotenv.Load(path); err != nil {
		logger.Warn("Error loading .env file - read system environment")
	}
}
