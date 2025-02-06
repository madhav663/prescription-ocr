package handlers_test

import (
	"bytes"
	//"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/madhav663/prescription-ocr/internal/api/handlers"
	"github.com/madhav663/prescription-ocr/internal/models"
)

func TestCreateMedication(t *testing.T) {
	
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()


	mock.ExpectQuery("INSERT INTO medications").
		WithArgs("Aspirin", "Ibuprofen", "Nausea").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	
	medicationModel := &models.MedicationModel{DB: db}

	
	handler := handlers.NewMedicationHandler(medicationModel)

	
	medication := &models.Medication{Name: "Aspirin", Alternatives: "Ibuprofen", SideEffects: "Nausea"}
	body, _ := json.Marshal(medication)
	req, err := http.NewRequest("POST", "/medications", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	
	rr := httptest.NewRecorder()
	handler.CreateMedication(rr, req)

	
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}


	var actual models.Medication
	if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	
	expected := models.Medication{ID: 1, Name: "Aspirin", Alternatives: "Ibuprofen", SideEffects: "Nausea"}

	
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %+v want %+v", actual, expected)
	}

	
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}