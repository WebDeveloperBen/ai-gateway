package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/insurgence-ai/llm-gateway/internal/provider"
)

type Core struct {
	MaxBody   int
	Transport http.RoundTripper
	Adapters  []provider.Adapter
}

func NewCoreWithAdapters(rt http.RoundTripper, adapters ...provider.Adapter) *Core {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return &Core{
		MaxBody:   1 << 20,
		Transport: rt,
		Adapters:  adapters,
	}
}

func (c *Core) StreamingHandler() func(ctx context.Context, _ *struct{}) (*huma.StreamResponse, error) {
	return func(ctx context.Context, _ *struct{}) (*huma.StreamResponse, error) {
		return &huma.StreamResponse{
			Body: func(hctx huma.Context) {
				// Build a placeholder request; Director will fully rewrite it.
				req, err := http.NewRequestWithContext(
					hctx.Context(),
					hctx.Method(),
					"http://placeholder", // avoids relying on hctx.URL().String()
					nil,
				)
				if err != nil {
					fail(hctx, http.StatusBadGateway, `{"title":"Bad Gateway","status":502,"detail":"build request failed"}`)
					return
				}

				hctx.EachHeader(func(n, v string) { req.Header.Add(n, v) })

				w := newHumaResponseWriter(hctx)
				if d, ok := hctx.BodyWriter().(interface{ SetWriteDeadline(time.Time) error }); ok {
					_ = d.SetWriteDeadline(time.Now().Add(60 * time.Second))
				}

				rp := &httputil.ReverseProxy{
					Transport:    c.Transport,
					Director:     c.makeDirector(hctx),
					ErrorHandler: writeProxyError,
				}
				rp.ServeHTTP(w, req)
			},
		}, nil
	}
}

func (c *Core) makeDirector(hctx huma.Context) func(*http.Request) {
	return func(req *http.Request) {
		// Use the real incoming path from Huma (not the placeholder).
		inURL := hctx.URL()
		path := inURL.Path
		req.URL.Path = path
		req.URL.RawQuery = inURL.RawQuery

		// 1) Try to find a matching adapter by Prefix() on segment boundaries.
		var (
			ad        provider.Adapter
			prefix    string
			prefixPos = -1
		)
		for _, a := range c.Adapters {
			pfx := a.Prefix()
			if pfx == "" || pfx == "/" {
				continue
			}
			if i := indexOfSegment(path, pfx); i >= 0 {
				ad, prefix, prefixPos = a, pfx, i
				break
			}
		}
		// Fallback: if exactly one adapter is registered, use it as default.
		if ad == nil && len(c.Adapters) == 1 {
			ad = c.Adapters[0]
			prefix = ""
			prefixPos = -1
		}
		if ad == nil {
			req.Header.Set("X-RP-Error", "no_adapter_for_path:"+path)
			req.Header.Set("X-RP-Adapters", strings.Join(listPrefixes(c.Adapters), ","))
			req.URL = mustParse("http://invalid/")
			return
		}

		// 2) Compute suffix starting at /v1/... AFTER the matched prefix.
		var tail string
		if prefixPos >= 0 {
			tail = path[prefixPos+len(prefix):]
			if !strings.HasPrefix(tail, "/") {
				tail = "/" + tail
			}
		} else {
			// default adapter case: tail is entire path
			tail = path
		}
		j := strings.Index(tail, "/v1/")
		if j < 0 {
			req.Header.Set("X-RP-Error", "no_v1_suffix after prefix "+prefix)
			req.URL = mustParse("http://invalid/")
			return
		}
		suffix := tail[j:] // "/v1/..."

		// 3) Snapshot body (so retries/RTs can reread) + build ReqInfo.
		var raw []byte
		if br := hctx.BodyReader(); br != nil {
			raw, _ = io.ReadAll(io.LimitReader(br, int64(c.MaxBody)))
		}
		// Attach body and be explicit about length & TE.
		req.Body = io.NopCloser(bytes.NewReader(raw))
		req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(raw)), nil }
		req.ContentLength = int64(len(raw))
		req.Header.Del("Transfer-Encoding")
		req.Header.Set("Content-Length", strconv.Itoa(len(raw)))
		// If you ever had an inbound Content-Encoding, drop it; we’re sending raw JSON upstream.
		req.Header.Del("Content-Encoding")

		info := provider.ReqInfo{
			Method: hctx.Method(),
			Path:   suffix,
			Query:  req.URL.RawQuery,
			Model:  extractModel(raw),
			Tenant: TenantFrom(hctx.Context()),
			App:    AppFrom(hctx.Context()),
		}

		// 4) Let the adapter rewrite to the real upstream.
		if err := ad.Rewrite(req, suffix, info); err != nil {
			req.Header.Set("X-RP-Error", "rewrite:"+escape(err.Error()))
			req.URL = mustParse("http://invalid/")
			return
		}

		// Ensure Host aligns with upstream host
		if req.URL.Host != "" {
			req.Host = req.URL.Host
		}
	}
}

// indexOfSegment finds `needle` inside `p` only when aligned on path segment boundaries.
// e.g. indexOfSegment("/api/azure/openai/v1/x", "/azure/openai") == 4.
func indexOfSegment(p, needle string) int {
	if needle == "" {
		return -1
	}
	i := strings.Index(p, needle)
	for i >= 0 {
		leftOK := (i == 0) || (p[i-1] == '/')
		right := i + len(needle)
		rightOK := (right == len(p)) || (p[right] == '/')
		if leftOK && rightOK {
			return i
		}
		next := strings.Index(p[i+1:], needle)
		if next < 0 {
			return -1
		}
		i += 1 + next
	}
	return -1
}

func listPrefixes(as []provider.Adapter) []string {
	out := make([]string, 0, len(as))
	for _, a := range as {
		out = append(out, a.Prefix())
	}
	return out
}

func writeProxyError(rw http.ResponseWriter, r *http.Request, err error) {
	rw.Header().Set("Content-Type", "application/problem+json")
	rw.WriteHeader(http.StatusBadGateway)
	detail := err.Error()
	if cause := r.Header.Get("X-RP-Error"); cause != "" {
		detail = cause + " | transport: " + detail
	}
	if have := r.Header.Get("X-RP-Adapters"); have != "" {
		detail += " | adapters: [" + have + "]"
	}
	_, _ = fmt.Fprintf(rw, `{"title":"Bad Gateway","status":502,"detail":%q}`, detail)
}

// ---- helpers ----

func mustParse(s string) *url.URL { u, _ := url.Parse(s); return u }

func extractModel(raw []byte) string {
	// Fast path for empty/very small bodies
	if len(raw) == 0 {
		return ""
	}
	// Try structured decode without allocating tons of stuff
	type m1 struct {
		Model string `json:"model"`
	}
	var a m1
	if json.Unmarshal(raw, &a) == nil && strings.TrimSpace(a.Model) != "" {
		return a.Model
	}
	// AOAI / legacy shapes sometimes use "deployment" or "engine"
	type m2 struct {
		Deployment string `json:"deployment"`
		Engine     string `json:"engine"`
		Model      string `json:"model"`
	}
	var b m2
	if json.Unmarshal(raw, &b) == nil {
		if s := strings.TrimSpace(b.Model); s != "" {
			return s
		}
		if s := strings.TrimSpace(b.Deployment); s != "" {
			return s
		}
		if s := strings.TrimSpace(b.Engine); s != "" {
			return s
		}
	}
	// Fallback: shallow map lookup without walking the whole JSON tree
	var obj map[string]any
	if json.Unmarshal(raw, &obj) == nil {
		for _, k := range []string{"model", "deployment", "engine"} {
			if v, ok := obj[k]; ok {
				if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
					return s
				}
			}
		}
	}
	return ""
}

func fail(hctx huma.Context, code int, body string) {
	hctx.SetStatus(code)
	hctx.SetHeader("Content-Type", "application/problem+json")
	_, _ = hctx.BodyWriter().Write([]byte(body))
}

func escape(s string) string {
	b, _ := json.Marshal(s)
	if len(b) >= 2 {
		return string(b[1 : len(b)-1])
	}
	return s
}
