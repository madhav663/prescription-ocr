package handlers

import (
   // "encoding/json"
    "net/http"

    "github.com/madhav663/prescription-ocr/internal/models"
)

// GetPrescriptionsHandler retrieves all stored prescriptions from the database
func GetPrescriptionsHandler(w http.ResponseWriter, r *http.Request) {
    prescriptions, err := models.GetPrescriptions()
    if err != nil {
        http.Error(w, "Failed to retrieve prescriptions", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(prescriptions)
}
