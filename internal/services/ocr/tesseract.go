package ocr

import (
	"fmt"
	"os/exec"
)


func ProcessImage(imagePath string) (string, error) {
	cmd := exec.Command("tesseract", imagePath, "stdout")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to process image: %w", err)
	}
	return string(output), nil
}