package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/madhav663/prescription-ocr/internal/api"
	"github.com/madhav663/prescription-ocr/internal/database/schema"
	"github.com/madhav663/prescription-ocr/internal/models"
	"github.com/madhav663/prescription-ocr/internal/services/llama"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func main() {
	
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	 
	requiredEnvs := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE", "LLAMA_API_URL"}
	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			log.Fatalf("Missing required environment variable: %s", env)
		}
	}

	
	llamaClient := llama.NewClient(os.Getenv("LLAMA_API_URL"))

	
	dbConfig := schema.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     getDBPort(),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	
	db, err := schema.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()
	log.Println("Successfully connected to the database.")

	
	medicationModel := &models.MedicationModel{DB: db}
	router := api.SetupRouter(medicationModel, llamaClient)

	
	port := getServerPort()
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getDBPort() int {
	dbPortStr := os.Getenv("DB_PORT")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Printf("Invalid DB_PORT: %s, defaulting to 5432", dbPortStr)
		return 5432
	}
	return dbPort
}

func getServerPort() string {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		return "8080"
	}
	return port
}
