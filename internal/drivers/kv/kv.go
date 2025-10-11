// Package kv defines the interface for a key-value store backend.
// Supports basic operations (Get, Set, Del), key existence check, and scanning keys.
package kv

import (
	"context"
	"strings"
	"time"
)

type KvStore interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Del(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)

	// Incr atomically increments the value at key by 1 and returns the new value
	// If key doesn't exist, it's set to 0 before incrementing (returns 1)
	Incr(ctx context.Context, key string) (int64, error)

	// IncrBy atomically increments the value at key by amount and returns the new value
	// If key doesn't exist, it's set to 0 before incrementing (returns amount)
	IncrBy(ctx context.Context, key string, amount int64) (int64, error)

	// Expire sets a TTL on an existing key
	// Returns true if TTL was set, false if key doesn't exist
	Expire(ctx context.Context, key string, ttl time.Duration) (bool, error)

	ScanGetAll(ctx context.Context, pattern string, count int64) (map[string]string, error)
	ScanAll(ctx context.Context, pattern string, count int64) ([]string, error)

	Close(ctx context.Context) error
}

/* ---------- Constants (namespaces) ---------- */

const (
	KeyUserSession         = "user:session"
	KeyAPIToken            = "api:token"
	KeyCachePrefix         = "cache:"
	KeyModelRegistryPrefix = "modelreg:"
)

const sep = ":"

/* ---------- Generic keyspace builder ---------- */

type Keyspace struct {
	// Prefix must end with sep (":"), NewKeyspace ensures this.
	Prefix string
}

func NewKeyspace(prefix string) Keyspace {
	if !strings.HasSuffix(prefix, sep) {
		prefix += sep
	}
	return Keyspace{Prefix: prefix}
}

// Key builds an exact key: ks.Prefix + parts joined by ":".
func (ks Keyspace) Key(parts ...string) string {
	if len(parts) == 0 {
		return ks.Prefix[:len(ks.Prefix)-1] // avoid trailing ":" for naked keyspace
	}
	return ks.Prefix + strings.Join(parts, sep)
}

// PrefixOf returns a prefix (always ends with ":") for the provided parts.
func (ks Keyspace) PrefixOf(parts ...string) string {
	if len(parts) == 0 {
		return ks.Prefix
	}
	return ks.Key(parts...) + sep
}

// PatternAll matches every key in this keyspace.
func (ks Keyspace) PatternAll() string {
	return ks.Prefix + "*"
}

// Pattern builds a SCAN MATCH pattern for a prefix of parts (escapes glob chars).
func (ks Keyspace) Pattern(parts ...string) string {
	if len(parts) == 0 {
		return ks.PatternAll()
	}
	var b strings.Builder
	b.WriteString(ks.Prefix)
	for i, p := range parts {
		if i > 0 {
			b.WriteString(sep)
		}
		b.WriteString(globEscape(p))
	}
	b.WriteString(sep)
	b.WriteString("*")
	return b.String()
}

// Escape only for PATTERNS (SCAN MATCH), never for exact keys.
func globEscape(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '*', '?', '[', ']', '\\':
			b.WriteByte('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}

/* ---------- Concrete keyspaces ---------- */

var (
	KSModelReg = NewKeyspace(KeyModelRegistryPrefix)
	KSCache    = NewKeyspace(KeyCachePrefix)
	KSUser     = NewKeyspace("user:")
	KSAPI      = NewKeyspace("api:")
)

// Thin wrappers

func KeyModel(tenant, model string) string  { return KSModelReg.Key(tenant, model) }
func TenantPrefix(tenant string) string     { return KSModelReg.PrefixOf(tenant) }
func PatternAll() string                    { return KSModelReg.PatternAll() }
func PatternTenantAll(tenant string) string { return KSModelReg.Pattern(tenant) }
