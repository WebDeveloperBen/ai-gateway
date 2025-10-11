package auth

import "context"

// Context keys for storing auth information
type contextKey string

const (
	contextKeyKeyID    contextKey = "api_key_id"
	contextKeyOrgID    contextKey = "org_id"
	contextKeyAppID    contextKey = "app_id"
	contextKeyUserID   contextKey = "user_id"
	contextKeyProvider contextKey = "provider"
	contextKeyModel    contextKey = "model_name"
	contextKeyPolicies contextKey = "policies"
)

// KeyData contains authenticated key information
type KeyData struct {
	KeyID  string
	OrgID  string
	AppID  string
	UserID string
}

// WithKeyID adds API key ID to context
func WithKeyID(ctx context.Context, keyID string) context.Context {
	return context.WithValue(ctx, contextKeyKeyID, keyID)
}

// GetKeyID retrieves API key ID from context
func GetKeyID(ctx context.Context) string {
	if val := ctx.Value(contextKeyKeyID); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// WithOrgID adds organization ID to context
func WithOrgID(ctx context.Context, orgID string) context.Context {
	return context.WithValue(ctx, contextKeyOrgID, orgID)
}

// GetOrgID retrieves organization ID from context
func GetOrgID(ctx context.Context) string {
	if val := ctx.Value(contextKeyOrgID); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// WithAppID adds application ID to context
func WithAppID(ctx context.Context, appID string) context.Context {
	return context.WithValue(ctx, contextKeyAppID, appID)
}

// GetAppID retrieves application ID from context
func GetAppID(ctx context.Context) string {
	if val := ctx.Value(contextKeyAppID); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// WithUserID adds user ID to context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, contextKeyUserID, userID)
}

// GetUserID retrieves user ID from context
func GetUserID(ctx context.Context) string {
	if val := ctx.Value(contextKeyUserID); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// WithProvider adds provider name to context
func WithProvider(ctx context.Context, provider string) context.Context {
	return context.WithValue(ctx, contextKeyProvider, provider)
}

// GetProvider retrieves provider name from context
func GetProvider(ctx context.Context) string {
	if val := ctx.Value(contextKeyProvider); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// WithModelName adds model name to context
func WithModelName(ctx context.Context, modelName string) context.Context {
	return context.WithValue(ctx, contextKeyModel, modelName)
}

// GetModelName retrieves model name from context
func GetModelName(ctx context.Context) string {
	if val := ctx.Value(contextKeyModel); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// WithPolicies adds loaded policies to context
func WithPolicies(ctx context.Context, policies interface{}) context.Context {
	return context.WithValue(ctx, contextKeyPolicies, policies)
}

// GetPolicies retrieves loaded policies from context
func GetPolicies(ctx context.Context) interface{} {
	return ctx.Value(contextKeyPolicies)
}
