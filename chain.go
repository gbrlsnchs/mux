package mux

import "net/http"

// Chain creates a chain of all middlewares plus a handler
// to create a single handler linked to each other through a "next" handler.
func Chain(fns ...MiddlewareFunc) MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
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
