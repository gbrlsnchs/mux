package mux

import "context"

// Params is a helper to access URL parameters.
func Params(ctx context.Context, ctxKey interface{}) map[string]string {
	return ctx.Value(ctxKey).(map[string]string)
}
