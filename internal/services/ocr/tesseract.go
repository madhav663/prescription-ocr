package ocr

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)


func ProcessImage(imagePath string) (string, error) {
	
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %s", imagePath)
	}

	log.Printf("Processing image with Tesseract: %s", imagePath)

	cmd := exec.Command("tesseract", imagePath, "stdout")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to process image: %s, Error: %v", imagePath, err)
		return "", fmt.Errorf("failed to process image: %w", err)
	}

	log.Printf("Tesseract processing complete for: %s", imagePath)
	return string(output), nil
}
