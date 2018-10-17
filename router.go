package mux

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gbrlsnchs/radix"
)

// MiddlewareFunc is a middleware adapter.
type MiddlewareFunc func(next http.Handler) http.Handler

// Router is an HTTP multiplexer.
type Router struct {
	path        string
	methods     map[string]*radix.Tree
	parentFns   []MiddlewareFunc
	fns         []MiddlewareFunc
	ctxKey      interface{}
	debug       bool
	placeholder byte
}

// NewRouter creates a new HTTP router.
func NewRouter(path string, ctxKey interface{}) *Router {
	if ctxKey == nil {
		panic("mux: context key is nil") // panic early in order to prevent panicking during application runtime
	}
	rt := &Router{
		methods:     make(map[string]*radix.Tree, 9),
		placeholder: ':',
		ctxKey:      ctxKey,
	}
	rt.path = rt.resolvePath(path)
	return rt
}

// Handle sets an HTTP request handler.
func (rt *Router) Handle(method, path string, handler http.Handler) {
	rt.handle(method, path, handler)
}

// HandleFunc sets an HTTP request handler function.
func (rt *Router) HandleFunc(method, path string, handler http.HandlerFunc) {
	rt.handle(method, path, handler)
}

// Router creates a subrouter.
func (rt *Router) Router(path string) *Router {
	return &Router{
		path:        rt.path + path,
		methods:     rt.methods,
		parentFns:   rt.fns,
		ctxKey:      rt.ctxKey,
		debug:       rt.debug,
		placeholder: rt.placeholder,
	}
}

// ServeHTTP implements the http.Handler interface by finding an endpoint in the trie.
// If there are any parameters, they are set to the request's context.
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if rt.methods[r.Method] == nil {
		http.NotFound(w, r)
		return
	}

	tr := rt.methods[r.Method]
	if n, p := tr.Get(r.URL.Path); n != nil {
		// Stores parameters in the request's context.
		if len(p) > 0 {
			ctx := r.Context()
			r = r.WithContext(context.WithValue(ctx, rt.ctxKey, p))
		}
		if handler, ok := n.Value.(http.Handler); ok {
			handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

// SetDebug sets the debug flag.
func (rt *Router) SetDebug(debug bool) {
	rt.debug = debug
}

// SetPlaceholder sets a different placeholder to be used in dynamic paths.
func (rt *Router) SetPlaceholder(c byte) {
	rt.placeholder = c
	for k := range rt.methods {
		rt.methods[k].SetBoundaries(rt.placeholder, '/')
	}
}

func (rt *Router) String() string {
	var bd strings.Builder
	for k, v := range rt.methods {
		fmt.Fprintf(&bd, "\n%s%v", k, v)
	}
	return bd.String()
}

// Use set middleware functions to run before each request.
func (rt *Router) Use(fns ...MiddlewareFunc) {
	rt.fns = fns
}

func (rt *Router) handle(method, path string, handler http.Handler) {
	fpn := rt.resolvePath(path)
	chain := NewChain(append(rt.parentFns, rt.fns...)...)(handler)
	m := rt.methods[method]
	if m == nil {
		flag := 0
		if rt.debug {
			flag = radix.Tdebug
		}
		m = radix.New(flag)
		m.SetBoundaries(rt.placeholder, '/')
		rt.methods[method] = m
	}
	m.Add(fpn, chain)
	m.Sort(radix.PrioritySort)
}

func (rt *Router) resolvePath(path string) string {
	var bd strings.Builder
	bd.WriteString(rt.path)
	bd.WriteString(path)
	path = bd.String()
	if path == "" {
		return "/"
	}
	for len(path) > 1 && path[0] == '/' && path[1] == '/' {
		path = path[1:]
	}
	return path
}
