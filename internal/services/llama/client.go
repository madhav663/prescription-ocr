package llama

import (
	"bytes"
	"encoding/json"
	"net/http"
)


type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}


func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}


func (c *Client) AnalyzeText(text string) (map[string]interface{}, error) {
	requestBody, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/analyze", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}