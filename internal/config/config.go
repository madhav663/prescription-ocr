// File: internal/config/config.go

package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
	Server struct {
		Port string
	}
	OCR struct {
		TesseractPath string
	}
	LLaMA struct {
		APIURL string
		APIKey string
	}
}

func LoadConfig() (*Config, error) {
	
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: .env file not found: %v\n", err)
	}

	config := &Config{}

	
	config.Database.Host = getEnv("DB_HOST", "localhost")
	config.Database.Port = getEnv("DB_PORT", "5432")
	config.Database.User = getEnv("DB_USER", "postgres")
	config.Database.Password = getEnv("DB_PASSWORD", "")
	config.Database.DBName = getEnv("DB_NAME", "prescription_ocr")
	config.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	
	config.Server.Port = getEnv("SERVER_PORT", "8080")

	
	config.OCR.TesseractPath = getEnv("TESSERACT_PATH", "/usr/bin/tesseract")

	
	config.LLaMA.APIURL = getEnv("LLAMA_API_URL", "")
	config.LLaMA.APIKey = getEnv("LLAMA_API_KEY", "")

	
	if config.LLaMA.APIURL == "" {
		return nil, fmt.Errorf("LLAMA_API_URL is required")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}