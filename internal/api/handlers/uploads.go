package handlers

import (
	"encoding/json"
	
	"io"
	"log"
	"net/http"
	"os"

	"github.com/madhav663/prescription-ocr/internal/models"
	"github.com/madhav663/prescription-ocr/internal/services/ocr"
)

func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, `{"error": "Failed to retrieve image"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadPath := "uploads/sample_text.png"
	outFile, err := os.Create(uploadPath)
	if err != nil {
		http.Error(w, `{"error": "Failed to save image"}`, http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, `{"error": "Failed to copy image data"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("✅ File uploaded successfully: %s", uploadPath)

	extractedText, err := ocr.ProcessImage(uploadPath)
	if err != nil {
		http.Error(w, `{"error": "Failed to extract text"}`, http.StatusInternalServerError)
		return
	}

	// Save extracted text to DB
	prescription := models.Prescription{
		OriginalImage: uploadPath,
		ExtractedText: extractedText,
	}
	if err := models.SavePrescription(&prescription); err != nil {
		http.Error(w, `{"error": "Failed to store prescription"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":        "File uploaded and processed successfully",
		"extracted_text": extractedText,
		"prescription_id": prescription.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
