package gateway

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type RTFunc func(*http.Request) (*http.Response, error)

func (f RTFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func Chain(base http.RoundTripper, mws ...func(http.RoundTripper) http.RoundTripper) http.RoundTripper {
	rt := base
	for i := len(mws) - 1; i >= 0; i-- {
		rt = mws[i](rt)
	}
	return rt
}

type Authenticator interface {
	Authenticate(*http.Request) (tenant, app string, err error)
}

func WithAuth(a Authenticator) func(http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return RTFunc(func(r *http.Request) (*http.Response, error) {
			tenant, app, err := a.Authenticate(r)
			if err != nil {
				return deny(401, "unauthorized"), nil
			}
			ctx := context.WithValue(r.Context(), ctxTenantKey{}, tenant)
			ctx = context.WithValue(ctx, ctxAppKey{}, app)
			return next.RoundTrip(r.WithContext(ctx))
		})
	}
}

type Limiter interface {
	Allow(*http.Request) (retryAfter time.Duration, ok bool)
}

func WithRateLimit(l Limiter) func(http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return RTFunc(func(r *http.Request) (*http.Response, error) {
			if retry, ok := l.Allow(r); !ok {
				return denyWithRetry(429, "rate limit exceeded", retry), nil
			}
			return next.RoundTrip(r)
		})
	}
}

type Metrics interface {
	Record(req *http.Request, resp *http.Response, d time.Duration)
}

func WithMetrics(m Metrics) func(http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return RTFunc(func(r *http.Request) (*http.Response, error) {
			start := time.Now()
			resp, err := next.RoundTrip(r)
			if err == nil {
				m.Record(r, resp, time.Since(start))
			}
			return resp, err
		})
	}
}

func deny(code int, msg string) *http.Response {
	return &http.Response{
		StatusCode: code,
		Status:     fmt.Sprintf("%d %s", code, http.StatusText(code)),
		Header:     http.Header{"Content-Type": []string{"application/problem+json"}},
		Body:       io.NopCloser(strings.NewReader(fmt.Sprintf(`{"title":"%s","status":%d}`, msg, code))),
	}
}

func denyWithRetry(code int, msg string, retryAfter time.Duration) *http.Response {
	h := http.Header{"Content-Type": []string{"application/problem+json"}}
	if retryAfter > 0 {
		h.Set("Retry-After", fmt.Sprintf("%d", int(retryAfter.Seconds())))
	}
	return &http.Response{
		StatusCode: code,
		Status:     fmt.Sprintf("%d %s", code, http.StatusText(code)),
		Header:     h,
		Body:       io.NopCloser(strings.NewReader(fmt.Sprintf(`{"title":"%s","status":%d}`, msg, code))),
	}
}
