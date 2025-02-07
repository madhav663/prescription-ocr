package llama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type AnalysisRequest struct {
	Text string `json:"text"`
}

type AnalysisResponse struct {
	Analysis    map[string]interface{} `json:"analysis"`
	Confidence  float64                `json:"confidence"`
	ProcessedAt time.Time              `json:"processed_at"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) AnalyzeText(ctx context.Context, text string) (*AnalysisResponse, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	requestBody, err := json.Marshal(AnalysisRequest{Text: text})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/analyze", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result AnalysisResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.ProcessedAt.IsZero() {
		result.ProcessedAt = time.Now()
	}

	return &result, nil
}

// ✅ **Fix: Add `AnalyzeMedication` to `client.go`**
func (c *Client) AnalyzeMedication(ctx context.Context, medicationName string) (map[string]interface{}, error) {
	// Structured prompt for medication analysis
	prompt := fmt.Sprintf(`Analyze the following medication:
	Name: %s
	
	Provide:
	1. Common uses
	2. Typical dosage
	3. Side effects
	4. Interactions
	5. Precautions`, medicationName)

	response, err := c.AnalyzeText(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze medication: %w", err)
	}

	return map[string]interface{}{
		"medication_name": medicationName,
		"analysis":        response.Analysis,
		"confidence":      response.Confidence,
		"analyzed_at":     response.ProcessedAt,
	}, nil
}
