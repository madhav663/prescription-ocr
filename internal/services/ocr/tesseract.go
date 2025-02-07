package ocr

import (
	"fmt"
	"log"
	"os/exec"
)

// ProcessImage extracts text using OCR
func ProcessImage(imagePath string) (string, error) {
	log.Println("🛠️ Running Tesseract on:", imagePath)

	cmd := exec.Command("tesseract", imagePath, "stdout", "--psm", "6", "--oem", "3")
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf(" Tesseract Error: %s\n", string(output))
		return "", fmt.Errorf("failed to process image: %s", string(output))
	}

	log.Println("✅ Extracted Text:", string(output))
	return string(output), nil
}
