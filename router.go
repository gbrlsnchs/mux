package mux

import (
	"net/http"
)

// Router is a prefix router created from an already existent mux.
type Router struct {
	path []byte
	m    *Mux
}

// Handle handles an HTTP handler according to a method and a path.
func (rt *Router) Handle(method, path string, handler http.Handler) {
	rt.m.handle(method, append(rt.path, path...), handler)
}
