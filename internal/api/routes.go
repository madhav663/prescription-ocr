package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/madhav663/prescription-ocr/internal/api/handlers"
	"github.com/madhav663/prescription-ocr/internal/services/ocr"
)

// SetupRouter initializes API routes
func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// ✅ API route for image upload
	mux.HandleFunc("/upload", UploadImageHandler)

	return mux
}


func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	
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

	// ✅ Save the file temporarily
	savePath := filepath.Join("uploads", header.Filename)
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"text": %q}`, ocrResult)))
}
