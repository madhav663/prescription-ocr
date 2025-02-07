package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/madhav663/prescription-ocr/internal/api"
	"github.com/madhav663/prescription-ocr/internal/database/schema"
	"github.com/madhav663/prescription-ocr/internal/models"
	"github.com/madhav663/prescription-ocr/internal/services/llama"
	_ "github.com/lib/pq"
)


func connectWithRetry(dbConfig schema.DBConfig, retries int) (*sql.DB, error) {
	var db *sql.DB
	var err error
	for i := 0; i < retries; i++ {
		log.Printf("🛠️ Attempting to connect to DB (Attempt %d/%d)...", i+1, retries)
		db, err = schema.NewDatabase(dbConfig)
		if err == nil {
			log.Println("✅ Database connection successful.")
			return db, nil
		}
		log.Printf("⚠️ Database connection attempt %d failed: %v", i+1, err)
		time.Sleep(time.Second * time.Duration(i+1)) 
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
	if models.DB == nil {
		t.Fatal(" Database connection is nil! Make sure models.DB is initialized.")
	}

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
	_, _ = io.Copy(part, file)
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	router := api.SetupRouter(&models.MedicationModel{DB: models.DB}, llama.NewClient(os.Getenv("LLAMA_API_URL")))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	fmt.Println("Raw Response:", rr.Body.String())

	var jsonResponse map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &jsonResponse)
	if err != nil {
		t.Fatalf(" Failed to parse JSON response: %v\nRaw Response: %s", err, rr.Body.String())
	}

	if jsonResponse["extracted_text"] == "" {
		t.Errorf("Expected extracted text, but got empty response")
	}
}
