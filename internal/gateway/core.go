package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type Core struct {
	Router  Router
	Client  *http.Client
	MaxBody int
}

func NewCore(router Router, transport http.RoundTripper) *Core {
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &Core{
		Router:  router,
		Client:  &http.Client{Transport: transport}, // no global timeout (streaming)
		MaxBody: 1 << 20,
	}
}

func (c *Core) StreamingHandler() func(ctx context.Context, _ *struct{}) (*huma.StreamResponse, error) {
	return func(ctx context.Context, _ *struct{}) (*huma.StreamResponse, error) {
		return &huma.StreamResponse{
			Body: func(hctx huma.Context) {
				inURL := hctx.URL()

				raw, err := io.ReadAll(io.LimitReader(hctx.BodyReader(), int64(c.MaxBody)))
				if err != nil {
					fail(hctx, http.StatusBadRequest, `{"title":"Bad Request","status":400,"detail":"failed to read body"}`)
					return
				}

				model := extractModel(raw)
				info := RequestInfo{
					Method: hctx.Method(),
					Path:   inURL.Path,
					Query:  inURL.RawQuery,
					Model:  model,
					Tenant: TenantFrom(hctx.Context()),
					App:    AppFrom(hctx.Context()),
				}

				target, mutate, err := c.Router.Route(info)
				if err != nil {
					fail(hctx, http.StatusBadRequest, `{"title":"Bad Request","status":400,"detail":"`+escape(err.Error())+`"}`)
					return
				}

				upReq, err := http.NewRequestWithContext(hctx.Context(), hctx.Method(), target, bytes.NewReader(raw))
				if err != nil {
					fail(hctx, http.StatusBadGateway, `{"title":"Bad Gateway","status":502,"detail":"failed to build upstream request"}`)
					return
				}

				hctx.EachHeader(func(n, v string) { upReq.Header.Add(n, v) })
				if mutate != nil {
					if err := mutate(hctx, upReq, raw); err != nil {
						fail(hctx, http.StatusBadRequest, `{"title":"Bad Request","status":400,"detail":"request not allowed"}`)
						return
					}
				}

				resp, err := c.Client.Do(upReq)
				if err != nil {
					fail(hctx, http.StatusBadGateway, `{"title":"Bad Gateway","status":502,"detail":"upstream request failed"}`)
					return
				}
				defer resp.Body.Close()

				hctx.SetStatus(resp.StatusCode)
				for k, vals := range resp.Header {
					for _, v := range vals {
						hctx.AppendHeader(k, v)
					}
				}

				w := hctx.BodyWriter()
				if d, ok := w.(interface{ SetWriteDeadline(time.Time) error }); ok {
					_ = d.SetWriteDeadline(time.Now().Add(60 * time.Second))
				}
				_, _ = io.Copy(w, resp.Body)
			},
		}, nil
	}
}

func extractModel(raw []byte) string {
	var s struct {
		Model string `json:"model"`
	}
	_ = json.Unmarshal(raw, &s)
	return s.Model
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
