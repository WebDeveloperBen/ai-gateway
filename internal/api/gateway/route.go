// Package provider provides the HTTP endpoints that proxy client requests
// through the gateway to upstream LLM providers.
package gateway

import (
	"net/http"

	"github.com/WebDeveloperBen/ai-gateway/internal/gateway"
	"github.com/danielgtaylor/huma/v2"
)

var openAIPaths = []string{
	"/v1/chat/completions",
	"/v1/completions",
	"/v1/embeddings",
	// TODO: add more as needed
}

func RegisterProvider(grp *huma.Group, base string, core *gateway.Core) {
	h := core.StreamingHandler()
	for _, p := range openAIPaths {
		huma.Register(grp, huma.Operation{
			OperationID:   "proxy-" + base + p,
			Method:        http.MethodPost,
			Path:          base + p,
			Summary:       "Proxy " + base + p,
			DefaultStatus: http.StatusOK,
			Tags:          []string{"Proxy", base},
		}, h)
	}
}
