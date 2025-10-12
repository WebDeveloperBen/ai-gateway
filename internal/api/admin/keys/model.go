package keys

import (
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

type APIKey struct {
	KeyID      string
	Tenant     string
	App        string
	Status     model.KeyStatus
	ExpiresAt  *time.Time
	LastUsedAt *time.Time
	LastFour   string
	Metadata   map[string]any
	CreatedAt  time.Time
}

// Request/Response body types
type MintKeyRequestBody struct {
	OrgID    string
	AppID    string
	UserID   string
	TTL      time.Duration
	Prefix   string
	Metadata map[string]any
}

type MintKeyResponseBody struct {
	Token string `json:"token"`
	Key   APIKey `json:"key"`
}

// Huma request/response types
type MintKeyRequest struct {
	Body MintKeyRequestBody `json:"body"`
}

type MintKeyResponse struct {
	Body MintKeyResponseBody `json:"body"`
}

type RevokeKeyRequest struct {
	KeyID string `path:"key_id" required:"true"`
}
