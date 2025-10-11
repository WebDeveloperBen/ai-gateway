package model

import (
	"time"

	"github.com/google/uuid"
)

type UsageMetric struct {
	ID                uuid.UUID
	OrgID             uuid.UUID
	AppID             uuid.UUID
	APIKeyID          uuid.UUID
	ModelID           *uuid.UUID
	Provider          string
	ModelName         string
	PromptTokens      int
	CompletionTokens  int
	TotalTokens       int
	RequestSizeBytes  int
	ResponseSizeBytes int
	Timestamp         time.Time
}

type TokenUsage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}
