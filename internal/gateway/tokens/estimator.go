// Package tokens provides token estimation and counting for LLM requests
package tokens

import (
	"encoding/json"
	"fmt"

	lru "github.com/hashicorp/golang-lru/v2"
	tiktoken "github.com/pkoukk/tiktoken-go"
)

const (
	MessageOverheadTokens = 4
	ArrayOverheadTokens   = 3
)

// Estimator provides token estimation for different models
type Estimator struct {
	encodings *lru.Cache[string, *tiktoken.Tiktoken]
}

// NewEstimator creates a new token estimator with LRU cache
func NewEstimator() *Estimator {
	cache, _ := lru.New[string, *tiktoken.Tiktoken](50)
	return &Estimator{
		encodings: cache,
	}
}

// EstimateRequest estimates tokens for a request body
func (e *Estimator) EstimateRequest(model string, body []byte) (int, error) {
	// Parse request body to extract messages/prompt
	var req struct {
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages,omitempty"`
		Prompt string `json:"prompt,omitempty"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		return 0, fmt.Errorf("failed to parse request body: %w", err)
	}

	// Get encoding for model
	encoding, err := e.getEncoding(model)
	if err != nil {
		return 0, err
	}

	var totalTokens int

	// Estimate from messages (chat completion format)
	if len(req.Messages) > 0 {
		for _, msg := range req.Messages {
			tokens := encoding.Encode(msg.Content, nil, nil)
			totalTokens += len(tokens)
			totalTokens += MessageOverheadTokens
		}
		totalTokens += ArrayOverheadTokens
		return totalTokens, nil
	}

	// Estimate from prompt (completion format)
	if req.Prompt != "" {
		tokens := encoding.Encode(req.Prompt, nil, nil)
		return len(tokens), nil
	}

	return 0, fmt.Errorf("no messages or prompt found in request")
}

// getEncoding gets or creates a tiktoken encoding for the given model
func (e *Estimator) getEncoding(model string) (*tiktoken.Tiktoken, error) {
	// Check LRU cache
	if enc, ok := e.encodings.Get(model); ok {
		return enc, nil
	}

	// Map model to encoding
	encodingName := getEncodingForModel(model)

	// Create encoding
	enc, err := tiktoken.GetEncoding(encodingName)
	if err != nil {
		return nil, fmt.Errorf("failed to get encoding for model %s: %w", model, err)
	}

	// Cache it (LRU will evict oldest if full)
	e.encodings.Add(model, enc)
	return enc, nil
}

// getEncodingForModel maps model names to tiktoken encoding names
func getEncodingForModel(model string) string {
	// GPT-4 models
	if model == "gpt-4" || model == "gpt-4-32k" || model == "gpt-4-turbo" || model == "gpt-4-turbo-preview" {
		return "cl100k_base"
	}

	// GPT-4o models
	if model == "gpt-4o" || model == "gpt-4o-mini" {
		return "o200k_base"
	}

	// GPT-3.5 models
	if model == "gpt-3.5-turbo" || model == "gpt-3.5-turbo-16k" {
		return "cl100k_base"
	}

	// Default to cl100k_base for unknown models
	return "cl100k_base"
}
