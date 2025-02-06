package api

import (
	"github.com/madhav663/prescription-ocr/internal/api/handlers"
	"github.com/madhav663/prescription-ocr/internal/models"
	"net/http"
)

func SetupRouter(model *models.MedicationModel) http.Handler {

	medicationHandler := handlers.NewMedicationHandler(model)

	mux := http.NewServeMux()

	mux.HandleFunc("/medications", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			medicationHandler.CreateMedication(w, r)
		case http.MethodGet:
			medicationHandler.GetMedication(w, r)
		case http.MethodPut:
			medicationHandler.UpdateMedication(w, r)
		case http.MethodDelete:
			medicationHandler.DeleteMedication(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
