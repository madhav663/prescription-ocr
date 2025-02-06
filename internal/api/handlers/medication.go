package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/madhav663/prescription-ocr/internal/models"
)

type MedicationHandler struct {
	Model *models.MedicationModel
}

func NewMedicationHandler(model *models.MedicationModel) *MedicationHandler {
	return &MedicationHandler{Model: model}
}

func (h *MedicationHandler) CreateMedication(w http.ResponseWriter, r *http.Request) {
	var medication models.Medication
	if err := json.NewDecoder(r.Body).Decode(&medication); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.Model.AddMedication(&medication); err != nil {
		http.Error(w, "Failed to add medication", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(medication)
}

func (h *MedicationHandler) GetMedication(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
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

	json.NewEncoder(w).Encode(medication)
}

func (h *MedicationHandler) UpdateMedication(w http.ResponseWriter, r *http.Request) {
	var medication models.Medication
	if err := json.NewDecoder(r.Body).Decode(&medication); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.Model.UpdateMedication(&medication); err != nil {
		http.Error(w, "Failed to update medication", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(medication)
}

func (h *MedicationHandler) DeleteMedication(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid medication ID", http.StatusBadRequest)
		return
	}

	if err := h.Model.DeleteMedication(id); err != nil {
		http.Error(w, "Failed to delete medication", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
