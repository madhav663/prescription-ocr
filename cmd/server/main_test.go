package main

import (
    "database/sql"
    "fmt"
    "io"
    "log"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
    "time"

    _ "github.com/lib/pq"
    "github.com/madhav663/prescription-ocr/internal/api"
    "github.com/madhav663/prescription-ocr/internal/models"
    "github.com/madhav663/prescription-ocr/internal/services/llama"
)

func TestMainFunctionality(t *testing.T) {
    os.Setenv("DB_HOST", "127.0.0.1")
    os.Setenv("DB_PORT", "5432")
    os.Setenv("DB_USER", "testuser")       
    os.Setenv("DB_PASSWORD", "testpassword") 
    os.Setenv("DB_NAME", "testdb")       
    os.Setenv("DB_SSLMODE", "disable")

    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_SSLMODE"),
    )

    testDB, err := sql.Open("postgres", dsn)
    if err != nil {
        t.Fatalf("Failed to connect to test database: %v", err)
    }
    defer testDB.Close()

    if err := testDB.Ping(); err != nil {
        t.Fatalf("Database connection failed: %v", err)
    }
    log.Println("Connected to test database.")

    _, err = testDB.Exec(`
        CREATE TABLE IF NOT EXISTS medications (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL
        )
    `)
    if err != nil {
        t.Fatalf("Failed to create test table: %v", err)
    }

    var insertedID int
    tx, err := testDB.Begin()
    if err != nil {
        t.Fatalf("Failed to start transaction: %v", err)
    }

    err = tx.QueryRow("INSERT INTO medications (name) VALUES ('TestMedication') RETURNING id").Scan(&insertedID)
    if err != nil {
        tx.Rollback()
        t.Fatalf("Failed to insert test data: %v", err)
    }

    if err := tx.Commit(); err != nil {
        t.Fatalf("Failed to commit transaction: %v", err)
    }
    log.Printf("Inserted test medication with ID: %d", insertedID)

    medicationModel := &models.MedicationModel{DB: testDB}
    llamaClient := llama.NewClient(os.Getenv("LLAMA_API_URL"))
    router := api.SetupRouter(medicationModel, llamaClient)

    server := httptest.NewServer(router)
    defer server.Close()

    log.Printf("Waiting for database commit...")
    time.Sleep(2 * time.Second)

    requestURL := fmt.Sprintf("%s/medications?id=%d", server.URL, insertedID)
    log.Printf("Sending request to: %s", requestURL)

    resp, err := http.Get(requestURL)
    if err != nil {
        t.Fatalf("Failed to make request: %v", err)
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    log.Printf("Response body: %s", body)

    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status OK, got %v", resp.StatusCode)
    }
}

