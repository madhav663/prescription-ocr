package main

import (
	"log"
	"net/http"
	"os"

	"github.com/madhav663/prescription-ocr/internal/api"
	"github.com/madhav663/prescription-ocr/internal/database/schema"
	"github.com/madhav663/prescription-ocr/internal/models"
	"github.com/madhav663/prescription-ocr/internal/services/llama"
	_ "github.com/lib/pq"
)

func main() {
	
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
		log.Fatalf(" Failed to connect to the database: %v", err)
	}
	defer db.Close()

	
	models.DB = db 

	
	medicationModel := &models.MedicationModel{DB: db}

	
	llamaClient := llama.NewClient(os.Getenv("LLAMA_API_URL"))

	
	router := api.SetupRouter(medicationModel, llamaClient)

	
	port := getServerPort()
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}


func getDBPort() int {
	port := os.Getenv("DB_PORT")
	if port == "" {
		return 5432 // Default PostgreSQL port
	}
	return 5432
}


func getServerPort() string {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		return "8080"
	}
	return port
}
