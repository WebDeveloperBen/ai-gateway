package model

import (
	"time"

	"github.com/google/uuid"
)

type PolicyType string

const (
	// Predefined policy types
	PolicyTypeRateLimit      PolicyType = "rate_limit"
	PolicyTypeTokenLimit     PolicyType = "token_limit"
	PolicyTypeModelAllowlist PolicyType = "model_allowlist"
	PolicyTypeRequestSize    PolicyType = "request_size"

	// Custom CEL policy
	PolicyTypeCustomCEL PolicyType = "custom_cel"
)

type Policy struct {
	ID         uuid.UUID
	OrgID      uuid.UUID
	AppID      uuid.UUID
	PolicyType PolicyType
	Config     map[string]any // JSON blob
	Enabled    bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// RateLimitConfig Policy-specific config structures
type RateLimitConfig struct {
	RequestsPerMinute int `json:"requests_per_minute"`
	TokensPerMinute   int `json:"tokens_per_minute"`
}

type TokenLimitConfig struct {
	MaxPromptTokens     int `json:"max_prompt_tokens"`
	MaxCompletionTokens int `json:"max_completion_tokens"`
	MaxTotalTokens      int `json:"max_total_tokens"`
}

type ModelAllowlistConfig struct {
	AllowedModelIDs []string `json:"allowed_model_ids"`
}

type RequestSizeConfig struct {
	MaxRequestBytes int `json:"max_request_bytes"`
}
