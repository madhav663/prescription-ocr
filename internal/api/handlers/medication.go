package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/madhav663/prescription-ocr/internal/models"
	"github.com/madhav663/prescription-ocr/internal/services/llama"
)



type MedicationHandler struct {
    Model       *models.MedicationModel
    LlamaClient *llama.Client
}


func NewMedicationHandler(model *models.MedicationModel, llamaClient *llama.Client) *MedicationHandler {
    return &MedicationHandler{
        Model:       model,
        LlamaClient: llamaClient,
    }
}


func (h *MedicationHandler) CreateMedication(w http.ResponseWriter, r *http.Request) {
    var medication models.Medication
    if err := json.NewDecoder(r.Body).Decode(&medication); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    
    ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
    defer cancel()

   
    if medication.Name != "" && h.LlamaClient != nil {
        analysis, err := h.LlamaClient.AnalyzeMedication(ctx, medication.Name)
        if err != nil {
           
            log.Printf("Failed to analyze medication: %v", err)
        } else {
            medication.Analysis = analysis
        }
    }

    if err := h.Model.AddMedication(&medication); err != nil {
        http.Error(w, "Failed to add medication", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(medication)
}



func (h *MedicationHandler) GetMedication(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    if idStr == "" {
        http.Error(w, "Medication ID is required", http.StatusBadRequest)
        return
    }

    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid medication ID", http.StatusBadRequest)
        return
    }

    medication, err := h.Model.GetMedication(id)
    if err != nil {
        http.Error(w, "Medication not found", http.StatusNotFound)
        return
    }

    if r.URL.Query().Get("analyze") == "true" && h.LlamaClient != nil {
        analysis, err := h.LlamaClient.AnalyzeMedication(r.Context(), medication.Name)
        if err == nil {
            medication.Analysis = analysis
            
            h.Model.UpdateMedication(medication)
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(medication)
}


func (h *MedicationHandler) UpdateMedication(w http.ResponseWriter, r *http.Request) {
    var medication models.Medication
    if err := json.NewDecoder(r.Body).Decode(&medication); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

  
    if medication.ID == 0 {
        http.Error(w, "Medication ID is required", http.StatusBadRequest)
        return
    }


    existing, err := h.Model.GetMedication(medication.ID)
    if err != nil {
        http.Error(w, "Medication not found", http.StatusNotFound)
        return
    }

    
    if medication.Name != existing.Name && h.LlamaClient != nil {
        analysis, err := h.LlamaClient.AnalyzeMedication(r.Context(), medication.Name)
        if err == nil {
            medication.Analysis = analysis
        }
    }

    if err := h.Model.UpdateMedication(&medication); err != nil {
        http.Error(w, "Failed to update medication", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(medication)
}


func (h *MedicationHandler) DeleteMedication(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    if idStr == "" {
        http.Error(w, "Medication ID is required", http.StatusBadRequest)
        return
    }

    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid medication ID", http.StatusBadRequest)
        return
    }


    if _, err := h.Model.GetMedication(id); err != nil {
        http.Error(w, "Medication not found", http.StatusNotFound)
        return
    }

    if err := h.Model.DeleteMedication(id); err != nil {
        http.Error(w, "Failed to delete medication", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}