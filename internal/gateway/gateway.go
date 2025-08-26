// Package gateway implements a reverse proxy for LLM traffic, providing
// tenant-aware auth, rate limiting, circuit breaking, retries and metrics.
package gateway

import (
	"fmt"
	"net/http"

	"github.com/insurgence-ai/llm-gateway/internal/gateway/auth"
	"github.com/insurgence-ai/llm-gateway/internal/gateway/loadbalancing"
	"github.com/insurgence-ai/llm-gateway/internal/provider"
	"github.com/insurgence-ai/llm-gateway/internal/provider/azureopenai"
	"github.com/insurgence-ai/llm-gateway/internal/provider/openai"
)

type Core struct {
	MaxBody       int
	Transport     http.RoundTripper
	Adapters      []provider.Adapter
	Authenticator auth.KeyAuthenticator
}

func NewCoreWithAdapters(rt http.RoundTripper, auth auth.KeyAuthenticator, adapters ...provider.Adapter) *Core {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return &Core{
		MaxBody:       1 << 20,
		Transport:     rt,
		Adapters:      adapters,
		Authenticator: auth,
	}
}

// NewCoreWithRegistry builds Core from a model registry (via cache+db)
// and dynamically wires up provider adapters (azure, openai, etc).
// Any future providers can be added easily in this registration step.
func NewCoreWithRegistry(rt http.RoundTripper, auth auth.KeyAuthenticator, reg *Registry) *Core {
	deployments, err := reg.All("modelreg:*")
	if err != nil {
		panic(fmt.Sprintf("failed to load registry: %v", err))
	}
	azureAdapter := azureopenai.BuildProvider(deployments, loadbalancing.NewRoundRobinSelector())
	openaiAdapter := openai.BuildProvider(deployments, loadbalancing.NewRoundRobinSelector())
	adapters := []provider.Adapter{}
	if azureAdapter != nil {
		adapters = append(adapters, azureAdapter)
	}
	if openaiAdapter != nil {
		adapters = append(adapters, openaiAdapter)
	}
	return NewCoreWithAdapters(rt, auth, adapters...)
}
