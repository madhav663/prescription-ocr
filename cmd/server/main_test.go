package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/madhav663/prescription-ocr/internal/api"
	"github.com/madhav663/prescription-ocr/internal/database/schema"
	"github.com/madhav663/prescription-ocr/internal/models"
	"github.com/madhav663/prescription-ocr/internal/services/llama"
	"github.com/joho/godotenv" 
	_ "github.com/lib/pq"
)


func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(" No .env file found, using environment variables")
	}
}


func connectWithRetry(dbConfig schema.DBConfig, retries int) (*sql.DB, error) {
	var db *sql.DB
	var err error
	for i := 0; i < retries; i++ {
		log.Printf("🛠️ Attempting to connect to DB (Attempt %d/%d)...", i+1, retries)
		db, err = schema.NewDatabase(dbConfig)
		if err == nil {
			log.Println(" Database connection successful.")
			return db, nil
		}
		log.Printf("⚠️ Database connection attempt %d failed: %v", i+1, err)
		time.Sleep(time.Second * time.Duration(i+1)) // Exponential backoff
	}
	return nil, fmt.Errorf("database not reachable after %d attempts: %v", retries, err)
}


func TestMain(m *testing.M) {
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
		log.Fatalf(" Failed to connect to the test database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	
	models.DB = db

	
	os.Exit(m.Run())
}


func TestUploadImageHandler(t *testing.T) {
	filePath := "uploads/sample_text.png"
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf(" Test image not found: %s", filePath)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", "sample_text.png")
	if err != nil {
		t.Fatalf(" Failed to create form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatalf(" Failed to copy file data: %v", err)
	}

	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	medicationModel := &models.MedicationModel{DB: models.DB}
	llamaClient := llama.NewClient(os.Getenv("LLAMA_API_URL"))

	router := api.SetupRouter(medicationModel, llamaClient)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf(" Expected status OK, got %d", rr.Code)
	}

	var jsonResponse map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &jsonResponse)
	if err != nil {
		t.Fatalf(" Failed to parse JSON response: %v", err)
	}

	if jsonResponse["extracted_text"] == "" {
		t.Errorf(" Expected extracted text, but got empty response")
	}
}
