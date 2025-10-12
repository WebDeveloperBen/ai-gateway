package gateway_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/gateway"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway/auth"
	"github.com/stretchr/testify/require"
)

// Mock authenticator for testing
type mockAuthenticator struct {
	shouldFail bool
	keyData    *auth.KeyData
}

func (m *mockAuthenticator) Authenticate(r *http.Request) (string, *auth.KeyData, error) {
	if m.shouldFail {
		return "", nil, errors.New("unauthorized")
	}
	return "test-key-id", m.keyData, nil
}

// Mock limiter for testing
type mockLimiter struct {
	shouldAllow bool
	retryAfter  time.Duration
}

func (m *mockLimiter) Allow(r *http.Request) (time.Duration, bool) {
	return m.retryAfter, m.shouldAllow
}

// Mock metrics for testing
type mockMetrics struct {
	recorded bool
}

func (m *mockMetrics) Record(req *http.Request, resp *http.Response, d time.Duration) {
	m.recorded = true
}

func TestRTFunc_RoundTrip(t *testing.T) {
	called := false
	rt := gateway.RTFunc(func(r *http.Request) (*http.Response, error) {
		called = true
		return &http.Response{StatusCode: 200}, nil
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := rt.RoundTrip(req)

	require.NoError(t, err)
	require.True(t, called)
	require.Equal(t, 200, resp.StatusCode)
}

func TestChain(t *testing.T) {
	middleware1Called := false
	middleware2Called := false

	mw1 := func(next http.RoundTripper) http.RoundTripper {
		return gateway.RTFunc(func(r *http.Request) (*http.Response, error) {
			middleware1Called = true
			return next.RoundTrip(r)
		})
	}

	mw2 := func(next http.RoundTripper) http.RoundTripper {
		return gateway.RTFunc(func(r *http.Request) (*http.Response, error) {
			middleware2Called = true
			return next.RoundTrip(r)
		})
	}

	base := gateway.RTFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200}, nil
	})

	chained := gateway.Chain(base, mw1, mw2)
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := chained.RoundTrip(req)

	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.True(t, middleware1Called)
	require.True(t, middleware2Called)
}

func TestWithAuth_Success(t *testing.T) {
	authenticator := &mockAuthenticator{
		keyData: &auth.KeyData{
			KeyID:  "test-key-id",
			OrgID:  "test-org",
			AppID:  "test-app",
			UserID: "test-user",
		},
	}

	base := gateway.RTFunc(func(r *http.Request) (*http.Response, error) {
		// Check that context was populated
		ctx := r.Context()
		require.Equal(t, "test-key-id", auth.GetKeyID(ctx))
		require.Equal(t, "test-org", auth.GetOrgID(ctx))
		require.Equal(t, "test-app", auth.GetAppID(ctx))
		require.Equal(t, "test-user", auth.GetUserID(ctx))
		return &http.Response{StatusCode: 200}, nil
	})

	middleware := gateway.WithAuth(authenticator)
	wrapped := middleware(base)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := wrapped.RoundTrip(req)

	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

func TestWithAuth_Failure(t *testing.T) {
	authenticator := &mockAuthenticator{shouldFail: true}

	base := gateway.RTFunc(func(r *http.Request) (*http.Response, error) {
		t.Fatal("base should not be called")
		return nil, nil
	})

	middleware := gateway.WithAuth(authenticator)
	wrapped := middleware(base)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := wrapped.RoundTrip(req)

	require.NoError(t, err)
	require.Equal(t, 401, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	require.Contains(t, string(body), "unauthorized")
}

func TestWithRateLimit_Allow(t *testing.T) {
	limiter := &mockLimiter{shouldAllow: true}

	base := gateway.RTFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200}, nil
	})

	middleware := gateway.WithRateLimit(limiter)
	wrapped := middleware(base)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := wrapped.RoundTrip(req)

	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

func TestWithRateLimit_Deny(t *testing.T) {
	limiter := &mockLimiter{
		shouldAllow: false,
		retryAfter:  30 * time.Second,
	}

	base := gateway.RTFunc(func(r *http.Request) (*http.Response, error) {
		t.Fatal("base should not be called")
		return nil, nil
	})

	middleware := gateway.WithRateLimit(limiter)
	wrapped := middleware(base)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := wrapped.RoundTrip(req)

	require.NoError(t, err)
	require.Equal(t, 429, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	require.Contains(t, string(body), "rate limit exceeded")

	require.Equal(t, "30", resp.Header.Get("Retry-After"))
}

func TestWithMetrics(t *testing.T) {
	metrics := &mockMetrics{}

	base := gateway.RTFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200}, nil
	})

	middleware := gateway.WithMetrics(metrics)
	wrapped := middleware(base)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := wrapped.RoundTrip(req)

	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.True(t, metrics.recorded)
}

func TestWithMetrics_Error(t *testing.T) {
	metrics := &mockMetrics{}

	base := gateway.RTFunc(func(r *http.Request) (*http.Response, error) {
		return nil, http.ErrServerClosed
	})

	middleware := gateway.WithMetrics(metrics)
	wrapped := middleware(base)

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := wrapped.RoundTrip(req)

	require.Error(t, err)
	require.Nil(t, resp)
	require.False(t, metrics.recorded) // Should not record on error
}
