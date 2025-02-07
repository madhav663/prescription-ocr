package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/madhav663/prescription-ocr/internal/services/ocr"
)


func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	// Restrict to POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to read uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	
	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadDir, 0755)
		if err != nil {
			http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
			return
		}
	}

	
	savePath := filepath.Join(uploadDir, header.Filename)
	outFile, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	
	ocrResult, err := ocr.ProcessImage(savePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("OCR processing failed: %v", err), http.StatusInternalServerError)
		return
	}

	
	_ = os.Remove(savePath)

	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"text": ocrResult})
}
