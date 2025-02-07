package api

import (
	"net/http"

	"github.com/madhav663/prescription-ocr/internal/api/handlers"
	"github.com/madhav663/prescription-ocr/internal/models"
	"github.com/madhav663/prescription-ocr/internal/services/llama"
)

func SetupRouter(model *models.MedicationModel, llamaClient *llama.Client) *http.ServeMux {
	mux := http.NewServeMux()

	medicationHandler := handlers.NewMedicationHandler(model, llamaClient)
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

	mux.HandleFunc("/upload", handlers.UploadImageHandler)
	mux.HandleFunc("/prescriptions", handlers.GetPrescriptionsHandler)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "API is running. Use /medications or /upload", http.StatusNotFound)
	})

	return mux
}


