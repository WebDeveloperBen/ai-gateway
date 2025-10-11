package tokens

import (
	"encoding/json"
	"fmt"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

// ProviderParser defines the interface for provider-specific response parsers
type ProviderParser interface {
	ParseResponse(body []byte) (*model.TokenUsage, error)
}

// OpenAIParser handles OpenAI and Azure OpenAI response format
type OpenAIParser struct{}

func (p *OpenAIParser) ParseResponse(body []byte) (*model.TokenUsage, error) {
	var resp struct {
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	if resp.Usage.TotalTokens == 0 {
		return nil, fmt.Errorf("no usage data in response")
	}

	return &model.TokenUsage{
		PromptTokens:     resp.Usage.PromptTokens,
		CompletionTokens: resp.Usage.CompletionTokens,
		TotalTokens:      resp.Usage.TotalTokens,
	}, nil
}

// AnthropicParser handles Anthropic Claude API response format
type AnthropicParser struct{}

func (p *AnthropicParser) ParseResponse(body []byte) (*model.TokenUsage, error) {
	var resp struct {
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse Anthropic response: %w", err)
	}

	totalTokens := resp.Usage.InputTokens + resp.Usage.OutputTokens
	if totalTokens == 0 {
		return nil, fmt.Errorf("no usage data in response")
	}

	return &model.TokenUsage{
		PromptTokens:     resp.Usage.InputTokens,
		CompletionTokens: resp.Usage.OutputTokens,
		TotalTokens:      totalTokens,
	}, nil
}

// CohereParser handles Cohere API response format
type CohereParser struct{}

func (p *CohereParser) ParseResponse(body []byte) (*model.TokenUsage, error) {
	var resp struct {
		Meta struct {
			BilledUnits struct {
				InputTokens  int `json:"input_tokens"`
				OutputTokens int `json:"output_tokens"`
			} `json:"billed_units"`
		} `json:"meta"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse Cohere response: %w", err)
	}

	totalTokens := resp.Meta.BilledUnits.InputTokens + resp.Meta.BilledUnits.OutputTokens
	if totalTokens == 0 {
		return nil, fmt.Errorf("no usage data in response")
	}

	return &model.TokenUsage{
		PromptTokens:     resp.Meta.BilledUnits.InputTokens,
		CompletionTokens: resp.Meta.BilledUnits.OutputTokens,
		TotalTokens:      totalTokens,
	}, nil
}

// GoogleParser handles Google Gemini API response format
type GoogleParser struct{}

func (p *GoogleParser) ParseResponse(body []byte) (*model.TokenUsage, error) {
	var resp struct {
		UsageMetadata struct {
			PromptTokenCount     int `json:"promptTokenCount"`
			CandidatesTokenCount int `json:"candidatesTokenCount"`
			TotalTokenCount      int `json:"totalTokenCount"`
		} `json:"usageMetadata"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse Google response: %w", err)
	}

	if resp.UsageMetadata.TotalTokenCount == 0 {
		return nil, fmt.Errorf("no usage data in response")
	}

	return &model.TokenUsage{
		PromptTokens:     resp.UsageMetadata.PromptTokenCount,
		CompletionTokens: resp.UsageMetadata.CandidatesTokenCount,
		TotalTokens:      resp.UsageMetadata.TotalTokenCount,
	}, nil
}
