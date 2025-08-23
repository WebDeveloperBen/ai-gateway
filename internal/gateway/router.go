package gateway

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type RequestInfo struct {
	Method string
	Path   string
	Query  string
	Model  string
	Tenant string // set by auth middleware
	App    string // set by auth middleware
}

type MutateFunc func(hctx huma.Context, req *http.Request, rawBody []byte) error

type Router interface {
	Route(info RequestInfo) (targetURL string, mutate MutateFunc, err error)
}
