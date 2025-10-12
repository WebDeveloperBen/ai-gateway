package middleware

import (
	"bytes"
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/auth"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/policies"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockPolicyEngine implements policies.Engine interface for testing
type mockPolicyEngine struct {
	loadPoliciesResult []policies.Policy
	loadPoliciesError  error
	loadCallCount      int
}

func (m *mockPolicyEngine) LoadPolicies(ctx context.Context, appID string) ([]policies.Policy, error) {
	m.loadCallCount++
	return m.loadPoliciesResult, m.loadPoliciesError
}

// mockPolicy implements policies.Policy interface for testing
type mockPolicy struct {
	preCheckError error
	policyType    model.PolicyType
	preCheckCount int
}

func (m *mockPolicy) Type() model.PolicyType {
	return m.policyType
}

func (m *mockPolicy) PreCheck(ctx context.Context, reqCtx *policies.PreRequestContext) error {
	m.preCheckCount++
	return m.preCheckError
}

func (m *mockPolicy) PostCheck(ctx context.Context, reqCtx *policies.PostRequestContext) {
	// Not used in these tests
}

func TestPolicyEnforcer_NewPolicyEnforcer(t *testing.T) {
	mockEngine := &mockPolicyEngine{}
	enforcer := NewPolicyEnforcer(mockEngine)
	assert.NotNil(t, enforcer)
	assert.Equal(t, mockEngine, enforcer.engine)
}

func TestPolicyEnforcer_Middleware(t *testing.T) {
	// Mock next RoundTripper
	next := &testMockRoundTripper{
		responseBody: []byte(`{"status": "ok"}`),
	}

	t.Run("FastPath_NoAppID", func(t *testing.T) {
		mockEngine := &mockPolicyEngine{}
		enforcer := NewPolicyEnforcer(mockEngine)
		middleware := enforcer.Middleware(next)

		req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte(`{"model": "gpt-4"}`)))

		// No auth context set - should take fast path
		resp, err := middleware.RoundTrip(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		// Should not call policy engine
		assert.Equal(t, 0, mockEngine.loadCallCount)
	})

	t.Run("PolicyLoadFailure", func(t *testing.T) {
		mockEngine := &mockPolicyEngine{
			loadPoliciesError: errors.New("database error"),
		}
		enforcer := NewPolicyEnforcer(mockEngine)
		middleware := enforcer.Middleware(next)

		req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte(`{"model": "gpt-4"}`)))

		// Set up auth context
		ctx := auth.WithOrgID(req.Context(), "org-123")
		ctx = auth.WithAppID(ctx, "app-456")
		ctx = auth.WithKeyID(ctx, "key-789")
		req = req.WithContext(ctx)

		resp, err := middleware.RoundTrip(req)
		require.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode)

		assert.Equal(t, 1, mockEngine.loadCallCount)
	})

	t.Run("MissingParsedRequest", func(t *testing.T) {
		mockEngine := &mockPolicyEngine{
			loadPoliciesResult: []policies.Policy{},
		}
		enforcer := NewPolicyEnforcer(mockEngine)
		middleware := enforcer.Middleware(next)

		req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte(`{"model": "gpt-4"}`)))

		// Set up auth context but no parsed request
		ctx := auth.WithOrgID(req.Context(), "org-123")
		ctx = auth.WithAppID(ctx, "app-456")
		ctx = auth.WithKeyID(ctx, "key-789")
		req = req.WithContext(ctx)

		resp, err := middleware.RoundTrip(req)
		require.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode)

		assert.Equal(t, 1, mockEngine.loadCallCount)
	})

	t.Run("PolicyPreCheckFailure", func(t *testing.T) {
		mockPolicy := &mockPolicy{
			preCheckError: errors.New("policy violation"),
			policyType:    model.PolicyTypeRateLimit,
		}
		mockEngine := &mockPolicyEngine{
			loadPoliciesResult: []policies.Policy{mockPolicy},
		}
		enforcer := NewPolicyEnforcer(mockEngine)
		middleware := enforcer.Middleware(next)

		req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte(`{"model": "gpt-4"}`)))

		// Set up auth context
		ctx := auth.WithOrgID(req.Context(), "org-123")
		ctx = auth.WithAppID(ctx, "app-456")
		ctx = auth.WithKeyID(ctx, "key-789")
		req = req.WithContext(ctx)

		// Set up parsed request
		parsedReq := &auth.ParsedRequest{
			Model:           "gpt-4",
			EstimatedTokens: 100,
			RequestSize:     50,
		}
		ctx = auth.WithParsedRequest(ctx, parsedReq)
		req = req.WithContext(ctx)

		resp, err := middleware.RoundTrip(req)
		require.NoError(t, err)
		assert.Equal(t, 429, resp.StatusCode)

		assert.Equal(t, 1, mockEngine.loadCallCount)
		assert.Equal(t, 1, mockPolicy.preCheckCount)
	})

	t.Run("SuccessfulPolicyEnforcement", func(t *testing.T) {
		mockPolicy := &mockPolicy{
			policyType: model.PolicyTypeRateLimit,
		}
		mockEngine := &mockPolicyEngine{
			loadPoliciesResult: []policies.Policy{mockPolicy},
		}
		enforcer := NewPolicyEnforcer(mockEngine)
		middleware := enforcer.Middleware(next)

		req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte(`{"model": "gpt-4"}`)))

		// Set up auth context
		ctx := auth.WithOrgID(req.Context(), "org-123")
		ctx = auth.WithAppID(ctx, "app-456")
		ctx = auth.WithKeyID(ctx, "key-789")
		req = req.WithContext(ctx)

		// Set up parsed request
		parsedReq := &auth.ParsedRequest{
			Model:           "gpt-4",
			EstimatedTokens: 100,
			RequestSize:     50,
		}
		ctx = auth.WithParsedRequest(ctx, parsedReq)
		req = req.WithContext(ctx)

		resp, err := middleware.RoundTrip(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		assert.Equal(t, 1, mockEngine.loadCallCount)
		assert.Equal(t, 1, mockPolicy.preCheckCount)
	})

	t.Run("MultiplePolicies_AllPass", func(t *testing.T) {
		mockPolicy1 := &mockPolicy{policyType: model.PolicyTypeRateLimit}
		mockPolicy2 := &mockPolicy{policyType: model.PolicyTypeTokenLimit}
		mockEngine := &mockPolicyEngine{
			loadPoliciesResult: []policies.Policy{mockPolicy1, mockPolicy2},
		}
		enforcer := NewPolicyEnforcer(mockEngine)
		middleware := enforcer.Middleware(next)

		req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte(`{"model": "gpt-4"}`)))

		// Set up auth context
		ctx := auth.WithOrgID(req.Context(), "org-123")
		ctx = auth.WithAppID(ctx, "app-456")
		ctx = auth.WithKeyID(ctx, "key-789")
		req = req.WithContext(ctx)

		// Set up parsed request
		parsedReq := &auth.ParsedRequest{
			Model:           "gpt-4",
			EstimatedTokens: 100,
			RequestSize:     50,
		}
		ctx = auth.WithParsedRequest(ctx, parsedReq)
		req = req.WithContext(ctx)

		resp, err := middleware.RoundTrip(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		assert.Equal(t, 1, mockEngine.loadCallCount)
		assert.Equal(t, 1, mockPolicy1.preCheckCount)
		assert.Equal(t, 1, mockPolicy2.preCheckCount)
	})

	t.Run("FirstPolicyFails_StopsProcessing", func(t *testing.T) {
		mockPolicy1 := &mockPolicy{
			preCheckError: errors.New("first policy fails"),
			policyType:    model.PolicyTypeRateLimit,
		}
		mockPolicy2 := &mockPolicy{
			policyType: model.PolicyTypeTokenLimit,
		}
		mockEngine := &mockPolicyEngine{
			loadPoliciesResult: []policies.Policy{mockPolicy1, mockPolicy2},
		}
		enforcer := NewPolicyEnforcer(mockEngine)
		middleware := enforcer.Middleware(next)

		req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte(`{"model": "gpt-4"}`)))

		// Set up auth context
		ctx := auth.WithOrgID(req.Context(), "org-123")
		ctx = auth.WithAppID(ctx, "app-456")
		ctx = auth.WithKeyID(ctx, "key-789")
		req = req.WithContext(ctx)

		// Set up parsed request
		parsedReq := &auth.ParsedRequest{
			Model:           "gpt-4",
			EstimatedTokens: 100,
			RequestSize:     50,
		}
		ctx = auth.WithParsedRequest(ctx, parsedReq)
		req = req.WithContext(ctx)

		resp, err := middleware.RoundTrip(req)
		require.NoError(t, err)
		assert.Equal(t, 429, resp.StatusCode)

		assert.Equal(t, 1, mockEngine.loadCallCount)
		assert.Equal(t, 1, mockPolicy1.preCheckCount)
		assert.Equal(t, 0, mockPolicy2.preCheckCount) // Should not be called
	})
}

func TestDeny(t *testing.T) {
	t.Run("Deny400", func(t *testing.T) {
		resp := deny(400, "bad request")
		assert.Equal(t, 400, resp.StatusCode)
		assert.Equal(t, "400 Bad Request", resp.Status)
		assert.Equal(t, "application/problem+json", resp.Header.Get("Content-Type"))

		body := make([]byte, 100)
		n, err := resp.Body.Read(body)
		require.NoError(t, err)
		assert.Contains(t, string(body[:n]), `"title":"bad request"`)
		assert.Contains(t, string(body[:n]), `"status":400`)
	})

	t.Run("Deny500", func(t *testing.T) {
		resp := deny(500, "internal error")
		assert.Equal(t, 500, resp.StatusCode)
		assert.Equal(t, "500 Internal Server Error", resp.Status)
		assert.Equal(t, "application/problem+json", resp.Header.Get("Content-Type"))

		body := make([]byte, 100)
		n, err := resp.Body.Read(body)
		require.NoError(t, err)
		assert.Contains(t, string(body[:n]), `"title":"internal error"`)
		assert.Contains(t, string(body[:n]), `"status":500`)
	})
}
