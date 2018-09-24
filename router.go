package mux

import (
	"context"
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
	fns         []MiddlewareFunc
	ctxKey      interface{}
	debug       bool
	placeholder byte
}

// New creates a new HTTP multiplexer.
func New(path string) *Router {
	return &Router{
		path:        path,
		methods:     make(map[string]*radix.Tree, 9),
		placeholder: ':',
	}
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
		path:    rt.path + path,
		methods: rt.methods,
		fns:     rt.fns,
		ctxKey:  rt.ctxKey,
		debug:   rt.debug,
	}
}

// ServeHTTP implements the http.Handler interface by finding an endpoint in the trie.
// If there are any parameters, they are set to the request's context.
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if rt.methods[r.Method] == nil {
		http.NotFound(w, r)
		return
	}

	trie := rt.methods[r.Method]
	if n, p := trie.Get(r.URL.Path); n != nil {
		// Stores parameters in the request's context.
		if len(p) > 0 && rt.ctxKey != nil {
			ctx := r.Context()
			r = r.WithContext(context.WithValue(ctx, rt.ctxKey, Params(p)))
		}
		if handler, ok := n.Value.(http.Handler); ok {
			handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

// SetCtxKey sets the key for accessing parameters in the request's context object.
// The idea is to make this key unique, maybe using enums or a UUID.
func (rt *Router) SetCtxKey(v interface{}) {
	rt.ctxKey = v
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
		bd.WriteString(k)
		bd.WriteByte('\n')
		bd.WriteString(v.String())
	}
	return bd.String()
}

// Use set middleware functions to run before each request.
func (rt *Router) Use(fns ...MiddlewareFunc) {
	rt.fns = fns
}

func (rt *Router) handle(method, path string, handler http.Handler) {
	fpn := rt.resolvePath(path)
	chain := Chain(rt.fns...)(handler)
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
