package gateway

import "context"

type (
	ctxTenantKey struct{}
	ctxAppKey    struct{}
)

func TenantFrom(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(ctxTenantKey{}).(string); ok {
		return v
	}
	return ""
}

func AppFrom(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(ctxAppKey{}).(string); ok {
		return v
	}
	return ""
}
