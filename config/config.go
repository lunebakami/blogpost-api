package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TursoURL   string
	TursoToken string
}

func LoadConfig() *Config {

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

  TursoURL := os.Getenv("TURSO_DATABASE_URL")
  TursoToken := os.Getenv("TURSO_TOKEN")

	if TursoURL == "" || TursoToken == "" {
    log.Fatalf("TURSO_DATABASE_URL and TURSO_TOKEN must be set")
	}

	return &Config{
		TursoURL:   os.Getenv("TURSO_DATABASE_URL"),
		TursoToken: os.Getenv("TURSO_TOKEN"),
	}
}
