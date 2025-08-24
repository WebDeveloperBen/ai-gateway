package gateway

import (
	"net/http"
	"sync"

	"github.com/danielgtaylor/huma/v2"
)

// Map the standard golang http objects into a huma context

type humaResponseWriter struct {
	hctx    huma.Context
	written bool
	status  int
	headers http.Header
	lock    sync.Mutex
}

func newHumaResponseWriter(hctx huma.Context) *humaResponseWriter {
	return &humaResponseWriter{
		hctx:    hctx,
		headers: make(http.Header),
		status:  0,
	}
}

func (w *humaResponseWriter) Header() http.Header {
	return w.headers
}

func (w *humaResponseWriter) WriteHeader(statusCode int) {
	w.lock.Lock()

	defer w.lock.Unlock()

	if w.written {
		return
	}

	w.status = statusCode

	w.hctx.SetStatus(statusCode)

	for k, vals := range w.headers {
		for _, v := range vals {
			w.hctx.AppendHeader(k, v)
		}
	}
	w.written = true
}

func (w *humaResponseWriter) Write(p []byte) (int, error) {
	w.lock.Lock()
	if !w.written {
		w.WriteHeader(http.StatusOK)
	}

	w.lock.Unlock()

	return w.hctx.BodyWriter().Write(p)
}
