package provider

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

// EnsureAbsoluteBase turns things like "my-aoai" or "my-aoai.openai.azure.com"
// into "https://my-aoai.openai.azure.com". If full URL already, leave as-is.
func EnsureAbsoluteBase(base string, defaultHostSuffix string) (string, error) {
	b := strings.TrimRight(strings.TrimSpace(base), "/")
	if b == "" {
		return "", fmt.Errorf("empty base")
	}
	if strings.Contains(b, "://") {
		return b, nil
	}
	if strings.Contains(b, ".") {
		return "https://" + b, nil
	}
	if defaultHostSuffix != "" {
		return "https://" + b + "." + strings.TrimLeft(defaultHostSuffix, "."), nil
	}
	return "https://" + b, nil
}

// JoinURL builds an absolute URL by joining base + path segments, and merging query.
func JoinURL(base string, segs []string, q url.Values) (*url.URL, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	ps := []string{u.Path}
	for _, s := range segs {
		if s != "" {
			ps = append(ps, s)
		}
	}
	u.Path = path.Join(ps...)
	if q != nil {
		// merge with existing
		ex := u.Query()
		for k, vs := range q {
			for _, v := range vs {
				ex.Add(k, v)
			}
		}
		u.RawQuery = ex.Encode()
	}
	return u, nil
}

// SetUpstreamURL copies URL fields onto req (and aligns Host header).
func SetUpstreamURL(req *http.Request, u *url.URL) {
	req.URL.Scheme = u.Scheme
	req.URL.Host = u.Host
	req.URL.Path = u.Path
	req.URL.RawQuery = u.RawQuery
	req.Host = u.Host
}

// CopyQuery merges req.URL.RawQuery into the provided values.
func CopyQuery(req *http.Request) url.Values {
	if req.URL == nil {
		return url.Values{}
	}
	q := url.Values{}
	for k, vs := range req.URL.Query() {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	return q
}

// StripCallerAuth removes caller Authorization before upstream auth is applied.
func StripCallerAuth(h http.Header) {
	h.Del("Authorization")
}

// SetAPIKey sets an API key header if non-empty.
func SetAPIKey(h http.Header, headerName, key string) {
	if key != "" {
		h.Set(headerName, key)
	}
}

// ForceContentLength makes req explicit (helps picky origins).
func ForceContentLength(req *http.Request, bodyLen int) {
	req.ContentLength = int64(bodyLen)
	req.Header.Del("Transfer-Encoding")
	req.Header.Set("Content-Length", strconv.Itoa(bodyLen))
	req.Header.Del("Content-Encoding") // ensure we send raw JSON
}

// ModelOrDefault picks model, falling back to single/default entry.
// Returns chosen key and ok = true if something usable exists.
func ModelOrDefault(model string, hasExact func(string) bool, single func() (string, bool), fallbackExists bool, fallbackKey string) (string, bool) {
	m := strings.ToLower(strings.TrimSpace(model))
	if m != "" && hasExact(m) {
		return m, true
	}
	if k, ok := single(); ok {
		return k, true
	}
	if fallbackExists {
		return fallbackKey, true
	}
	return "", false
}
