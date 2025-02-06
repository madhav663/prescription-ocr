package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Madhav663/prescription-ocr/internal/api/handlers"
	"github.com/Madhav663/prescription-ocr/internal/api/middleware"
	"github.com/Madhav663/prescription-ocr/internal/api/routes"
	"github.com/Madhav663/prescription-ocr/internal/database/schema"
	"github.com/Madhav663/prescription-ocr/internal/models"
	"github.com/Madhav663/prescription-ocr/internal/services/llama"
	"github.com/Madhav663/prescription-ocr/internal/services/ocr"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		dbPort = 5432 // Default port
	}

	dbConfig := schema.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432,
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

	medicationModel := &models.MedicationModel{DB: db}
	ocrService := &ocr.Service{TesseractPath: os.Getenv("TESSERACT_PATH")}
	llamaClient := llama.NewClient(os.Getenv("LLAMA_API_URL"))

	medicationHandler := handlers.NewMedicationHandler(medicationModel, ocrService, llamaClient)
	router := routes.NewRouter(medicationHandler)
	handlerWithMiddleware := middleware.ApplyCORS(router)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, handlerWithMiddleware); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	fmt.Printf("everything is ok")
}
