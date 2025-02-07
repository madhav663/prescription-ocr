package llama

import (
	"context"
	"testing"
	"time"
)


type MockLlamaClient struct{}


func (m *MockLlamaClient) MockAnalyzeText(ctx context.Context, text string) (*AnalysisResponse, error) {
	time.Sleep(50 * time.Millisecond) 

	return &AnalysisResponse{
		Analysis: map[string]interface{}{
			"summary":  "Mock analysis of: " + text,
			"keywords": []string{"mock", "test", "analysis"},
		},
		Confidence:  0.99,
		ProcessedAt: time.Now(),
	}, nil
}


func (m *MockLlamaClient) MockAnalyzeMedication(ctx context.Context, medicationName string) (map[string]interface{}, error) {
	time.Sleep(50 * time.Millisecond) // Simulate API delay

	return map[string]interface{}{
		"medication_name": medicationName,
		"analysis": map[string]interface{}{
			"uses":         "Pain relief",
			"dosage":       "500mg once daily",
			"side_effects": "Drowsiness, dizziness",
			"interactions": "Avoid alcohol",
			"precautions":  "Not recommended for pregnant women",
		},
		"confidence":  0.99,
		"analyzed_at": time.Now(),
	}, nil
}


func TestMockAnalyzeText(t *testing.T) {
	mockClient := &MockLlamaClient{}
	ctx := context.Background()

	response, err := mockClient.MockAnalyzeText(ctx, "This is a test sentence.")
	if err != nil {
		t.Fatalf("MockAnalyzeText failed: %v", err)
	}

	if response.Confidence < 0.5 {
		t.Fatalf("MockAnalyzeText returned low confidence: %f", response.Confidence)
	}

	if len(response.Analysis) == 0 {
		t.Fatalf("MockAnalyzeText returned empty analysis")
	}
}


func TestMockAnalyzeMedication(t *testing.T) {
	mockClient := &MockLlamaClient{}
	ctx := context.Background()

	response, err := mockClient.MockAnalyzeMedication(ctx, "Paracetamol")
	if err != nil {
		t.Fatalf("MockAnalyzeMedication failed: %v", err)
	}

	if response["confidence"].(float64) < 0.5 {
		t.Fatalf("MockAnalyzeMedication returned low confidence: %f", response["confidence"].(float64))
	}

	if response["medication_name"] != "Paracetamol" {
		t.Fatalf("Expected medication_name to be 'Paracetamol', got %s", response["medication_name"])
	}
}

