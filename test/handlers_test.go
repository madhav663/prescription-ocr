package handlers_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "reflect"  // Add this for deep comparison

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/madhav663/prescription-ocr/internal/api/handlers"
    "github.com/madhav663/prescription-ocr/internal/models"
    "github.com/madhav663/prescription-ocr/internal/services/llama"
)

func TestCreateMedication(t *testing.T) {
    // Create a new mock database connection
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    // Set up database mock expectations
    mock.ExpectQuery("INSERT INTO medications").
        WithArgs("Aspirin", "Ibuprofen", "Nausea").
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

    // Create medication model with mock db
    medicationModel := &models.MedicationModel{DB: db}

    // Create mock Llama client - Fix for first error
    mockLlamaClient := llama.NewClient("http://mock-llama-url")

    // Create handler with both dependencies - Fix for first error
    handler := handlers.NewMedicationHandler(medicationModel, mockLlamaClient)

    // Create test medication
    medication := &models.Medication{
        Name:         "Aspirin",
        Alternatives: "Ibuprofen",
        SideEffects:  "Nausea",
    }
    body, _ := json.Marshal(medication)

    // Create request
    req, err := http.NewRequest("POST", "/medications", bytes.NewBuffer(body))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Create response recorder
    rr := httptest.NewRecorder()

    // Call the handler
    handler.CreateMedication(rr, req)

    // Check status code
    if status := rr.Code; status != http.StatusCreated {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
    }

    // Decode response
    var actual models.Medication
    if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
        t.Fatalf("could not decode response: %v", err)
    }

    // Create expected result
    expected := models.Medication{
        ID:           1,
        Name:         "Aspirin",
        Alternatives: "Ibuprofen",
        SideEffects:  "Nausea",
    }

    // Fix for second error: Compare fields individually instead of comparing entire structs
    if actual.ID != expected.ID {
        t.Errorf("handler returned unexpected ID: got %v want %v", actual.ID, expected.ID)
    }
    if actual.Name != expected.Name {
        t.Errorf("handler returned unexpected Name: got %v want %v", actual.Name, expected.Name)
    }
    if actual.Alternatives != expected.Alternatives {
        t.Errorf("handler returned unexpected Alternatives: got %v want %v", actual.Alternatives, expected.Alternatives)
    }
    if actual.SideEffects != expected.SideEffects {
        t.Errorf("handler returned unexpected SideEffects: got %v want %v", actual.SideEffects, expected.SideEffects)
    }

    // If you need to compare the Analysis field (which is a map), use reflect.DeepEqual
    if !reflect.DeepEqual(actual.Analysis, expected.Analysis) {
        t.Errorf("handler returned unexpected Analysis: got %v want %v", actual.Analysis, expected.Analysis)
    }

    // Verify all expectations were met
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectations: %s", err)
    }
}