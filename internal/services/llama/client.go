package llama

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// Client represents a Llama API client
type Client struct {
    BaseURL    string
    HTTPClient *http.Client
}

// AnalysisRequest represents the request structure for text analysis
type AnalysisRequest struct {
    Text string `json:"text"`
}

// AnalysisResponse represents the structured response from Llama
type AnalysisResponse struct {
    Analysis    map[string]interface{} `json:"analysis"`
    Confidence  float64               `json:"confidence"`
    ProcessedAt time.Time             `json:"processed_at"`
}

// NewClient creates a new Llama client with configured timeout
func NewClient(baseURL string) *Client {
    return &Client{
        BaseURL: baseURL,
        HTTPClient: &http.Client{
            Timeout: 30 * time.Second, // Add timeout
        },
    }
}

// AnalyzeText sends text to Llama API for analysis with context and improved error handling
func (c *Client) AnalyzeText(ctx context.Context, text string) (*AnalysisResponse, error) {
    // Validate input
    if text == "" {
        return nil, fmt.Errorf("text cannot be empty")
    }

    // Create request with context
    requestBody, err := json.Marshal(AnalysisRequest{Text: text})
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }

    // Create new request with context
    req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/analyze", bytes.NewBuffer(requestBody))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    // Set headers
    req.Header.Set("Content-Type", "application/json")

    // Send request
    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    // Check status code
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    // Parse response
    var result AnalysisResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    // Set processed time if not provided by API
    if result.ProcessedAt.IsZero() {
        result.ProcessedAt = time.Now()
    }

    return &result, nil
}

// AnalyzeMedication specifically analyzes medication text
func (c *Client) AnalyzeMedication(ctx context.Context, medicationName string) (map[string]interface{}, error) {
    // Create a structured prompt for medication analysis
    prompt := fmt.Sprintf(`Analyze the following medication:
    Name: %s
    
    Please provide:
    1. Common uses
    2. Typical dosage
    3. Side effects
    4. Interactions
    5. Precautions`, medicationName)

    // Call AnalyzeText with the structured prompt
    response, err := c.AnalyzeText(ctx, prompt)
    if err != nil {
        return nil, fmt.Errorf("failed to analyze medication: %w", err)
    }

    // Structure the medication-specific response
    result := map[string]interface{}{
        "medication_name": medicationName,
        "analysis":       response.Analysis,
        "confidence":     response.Confidence,
        "analyzed_at":    response.ProcessedAt,
    }

    return result, nil
}

// MockAnalyzeText provides a mock implementation for testing
func (c *Client) MockAnalyzeText(ctx context.Context, text string) (*AnalysisResponse, error) {
    // Simulate API delay
    time.Sleep(100 * time.Millisecond)

    return &AnalysisResponse{
        Analysis: map[string]interface{}{
            "summary": "Mock analysis of: " + text,
            "keywords": []string{"mock", "test", "analysis"},
        },
        Confidence:  0.95,
        ProcessedAt: time.Now(),
    }, nil
}