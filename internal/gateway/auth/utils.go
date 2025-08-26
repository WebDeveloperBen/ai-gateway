package auth

import (
	"net/http"
	"strings"

	"github.com/insurgence-ai/llm-gateway/internal/repository/keys"
)

func getHeaderToken(r *http.Request) string {
	if v := r.Header.Get("Authorization"); v != "" {
		if after, ok := strings.CutPrefix(v, "Bearer "); ok {
			return strings.TrimSpace(after)
		}
	}
	if v := r.Header.Get("X-API-Key"); v != "" {
		return strings.TrimSpace(v)
	}
	return ""
}

func splitToken(tok string) (string, string) {
	left, right, ok := strings.Cut(tok, ".")
	if !ok || left == "" || right == "" {
		return "", ""
	}
	return left, right
}

// padWork burns comparable CPU on early rejects to avoid an existence oracle.
func padWork(h keys.Hasher) error {
	_, _ = h.Hash([]byte("timing-pad"))
	return nil
}
