package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/madhav663/prescription-ocr/internal/api"
	"github.com/madhav663/prescription-ocr/internal/database/schema"
	"github.com/madhav663/prescription-ocr/internal/models"
	"github.com/madhav663/prescription-ocr/internal/services/llama"
)

func connectWithRetry(dbConfig schema.DBConfig, retries int) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < retries; i++ {
		log.Printf(" Attempting to connect to DB (Attempt %d/%d)...", i+1, retries)
		db, err = schema.NewDatabase(dbConfig)
		if err == nil {
			log.Println(" Database connection successful.")
			return db, nil
		}
		log.Printf("⚠️ Database connection attempt %d failed: %v", i+1, err)
		time.Sleep(time.Second * time.Duration(i+1)) // Exponential backoff
	}

	return nil, fmt.Errorf(" Database not reachable after %d attempts: %v", retries, err)
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system environment variables")
	}

	requiredEnvs := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE", "LLAMA_API_URL"}
	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			log.Fatalf(" Missing required environment variable: %s", env)
		}
	}

	dbConfig := schema.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := connectWithRetry(dbConfig, 5)
	if err != nil {
		log.Fatalf(" Failed to connect to the database: %v", err)
	}
	defer db.Close()

	models.DB = db

	llamaClient := llama.NewClient(os.Getenv("LLAMA_API_URL"))
	medicationModel := &models.MedicationModel{DB: db}

	router := api.SetupRouter(medicationModel, llamaClient)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf(" Failed to start server: %v", err)
	}
}
