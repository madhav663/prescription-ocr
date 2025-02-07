package main

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/madhav663/prescription-ocr/internal/api/handlers"
)

func TestUploadImageHandler(t *testing.T) {
	dir, _ := os.Getwd()
	imagePath := filepath.Join(dir, "..", "..", "uploads", "sample_text.png")

	fmt.Println(" Checking test image path:", imagePath)

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		t.Fatalf(" Test image not found: %s. Please add a test image in the 'uploads/' directory.", imagePath)
	}

	file, err := os.Open(imagePath)
	if err != nil {
		t.Fatalf(" Failed to open test image: %v", err)
	}
	defer file.Close()

	
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", filepath.Base(imagePath))
	if err != nil {
		t.Fatalf(" Failed to create form file: %v", err)
	}

	
	imageData, _ := os.ReadFile(imagePath)
	part.Write(imageData)
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	handlers.UploadImageHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf(" Expected status OK, got %d", rec.Code)
	} else {
		fmt.Println(" Test Passed: Image uploaded successfully!")
	}
}
