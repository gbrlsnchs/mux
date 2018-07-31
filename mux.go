package mux

import (
	"bytes"
	"context"
	"net/http"

	"github.com/gbrlsnchs/mux/internal/radix"
)

// MiddlewareFunc is a middleware adapter.
type MiddlewareFunc func(next http.Handler) http.Handler

// Mux is an HTTP multiplexer.
type Mux struct {
	path    []byte
	methods map[string]*radix.Tree
	fns     []MiddlewareFunc
	ctxKey  interface{}
}

// New creates a new HTTP multiplexer.
func New(path string) *Mux {
	return &Mux{
		path:    []byte(path),
		methods: make(map[string]*radix.Tree, 9),
	}
}

// Handle handles an HTTP handler according to a method and a path.
func (m *Mux) Handle(method, path string, handler http.Handler) {
	fpn := m.resolvePath([]byte(path))
	chain := Chain(m.fns...)(handler)
	if m.methods[method] == nil {
		m.methods[method] = radix.New()
	}
	m.methods[method].Add(fpn, chain)
}

// ServeHTTP implements the http.Handler interface by finding an endpoint in the trie.
// If there are any parameters, they are set to the request's context.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.methods[r.Method] == nil {
		http.NotFound(w, r)
		return
	}

	trie := m.methods[r.Method]
	if n, params := trie.Get([]byte(r.URL.Path)); n != nil {
		// Stores parameters in the request's context.
		if len(params) > 0 && m.ctxKey != nil {
			ctx := r.Context()
			r = r.WithContext(context.WithValue(ctx, m.ctxKey, Params(params)))
		}
		n.Handler.ServeHTTP(w, r)
		return
	}
	http.NotFound(w, r)
}

// SetCtxKey sets the key for accessing parameters in the request's context object.
// The idea is to make this key unique, maybe using enums or a UUID.
func (m *Mux) SetCtxKey(v interface{}) {
	m.ctxKey = v
}

// Use set middleware functions to run before each request.
func (m *Mux) Use(fns ...MiddlewareFunc) {
	m.fns = fns
}

func (m *Mux) resolvePath(path []byte) []byte {
	var (
		buf bytes.Buffer
		b   []byte
		l   int
	)
	buf.Write(m.path)
	buf.Write(path)
	if b, l = buf.Bytes(), buf.Len(); l > 1 && b[0] == '/' && b[1] == '/' {
		return b[1:]
	}
	return b
}
