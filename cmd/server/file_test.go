package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestCheckFileAccess(t *testing.T) {
	// ✅ Get the project root directory instead of `cmd/server/`
	rootDir, _ := os.Getwd()
	imagePath := filepath.Join(rootDir, "..", "..", "uploads", "sample_text.png")

	fmt.Println("🔍 Checking file path:", imagePath)

	// ✅ Check if the file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		t.Fatalf("❌ Test image not found at: %s", imagePath)
	} else {
		fmt.Println("✅ Test image found at:", imagePath)
	}

	// ✅ Try opening the file
	file, err := os.Open(imagePath)
	if err != nil {
		t.Fatalf("❌ Failed to open file: %v", err)
	} else {
		fmt.Println("✅ Successfully opened:", imagePath)
	}
	defer file.Close()
}
