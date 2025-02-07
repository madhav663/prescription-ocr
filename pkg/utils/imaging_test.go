package utils

import (
	"image"
	"image/color"
	"os"
	"testing"
)


func TestSaveImage(t *testing.T) {
	
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			img.Set(x, y, color.White)
		}
	}

	testPath := "test_image.jpg"
	err := SaveImage(img, testPath)
	if err != nil {
		t.Fatalf("SaveImage failed: %v", err)
	}

	
	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		t.Fatalf("SaveImage did not create file: %s", testPath)
	}

	
	os.Remove(testPath)
}
