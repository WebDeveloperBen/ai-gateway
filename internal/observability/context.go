package observability

import "context"

type observabilityKey struct{}

// WithObservability adds observability to context
func WithObservability(ctx context.Context, obs *Observability) context.Context {
	return context.WithValue(ctx, observabilityKey{}, obs)
}

// FromContext retrieves observability from context
func FromContext(ctx context.Context) *Observability {
	if obs, ok := ctx.Value(observabilityKey{}).(*Observability); ok {
		return obs
	}
	return newNoopObservability()
}
