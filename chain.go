package mux

import "net/http"

// NewChain creates a middleware chain, which is a middleware itself.
func NewChain(fns ...MiddlewareFunc) MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		if handler == nil {
			panic("mux: handler is nil")
		}
		return build(handler, fns)
	}
}

func build(handler http.Handler, fns []MiddlewareFunc) http.Handler {
	if len(fns) == 0 {
		return handler
	}
	index := len(fns) - 1
	return build(fns[index](handler), fns[:index])
}
