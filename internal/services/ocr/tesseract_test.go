package ocr

import (
	"testing"
)


func TestProcessImage(t *testing.T) {
	
	testImagePath := "sample_text.png" 

	
	output, err := ProcessImage(testImagePath)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}

	
	if len(output) == 0 {
		t.Fatalf("ProcessImage returned empty result")
	}
}
