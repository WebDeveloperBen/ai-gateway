// Package proxy provides the HTTP endpoints that proxy client requests
// through the gateway to upstream LLM providers.
package proxy

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/insurgence-ai/llm-gateway/internal/gateway"
)

func RegisterRoutes(grp *huma.Group, core *gateway.Core) {
	handler := core.StreamingHandler()

	huma.Register(grp, huma.Operation{
		OperationID:   "proxy-chat-completions",
		Method:        http.MethodPost,
		Path:          "/v1/chat/completions",
		Summary:       "Proxy Chat Completions",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Proxy"},
	}, handler)

	huma.Register(grp, huma.Operation{
		OperationID:   "proxy-completions",
		Method:        http.MethodPost,
		Path:          "/v1/completions",
		Summary:       "Proxy Completions",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Proxy"},
	}, handler)

	huma.Register(grp, huma.Operation{
		OperationID:   "embeddings",
		Method:        http.MethodPost,
		Path:          "/v1/embeddings",
		Summary:       "Proxy Embeddings",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Proxy"},
	}, handler)

	// add more /v1/* here as you support them
}
