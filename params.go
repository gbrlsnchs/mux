package mux

import "context"

// Params is a helper to access URL parameters.
func Params(ctx context.Context, ctxKey interface{}) map[string]string {
	if p, ok := ctx.Value(ctxKey).(map[string]string); ok {
		return p
	}
	return map[string]string{}
}
