// config/config.go
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURI string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		return nil, fmt.Errorf("DB_URI environment variable is not set")
	}

	return &Config{
		DBURI: dbURI,
	}, nil
}
