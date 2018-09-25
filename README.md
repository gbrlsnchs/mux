# mux (HTTP multiplexer for Go)
[![Build Status](https://travis-ci.org/gbrlsnchs/mux.svg?branch=master)](https://travis-ci.org/gbrlsnchs/mux)
[![Sourcegraph](https://sourcegraph.com/github.com/gbrlsnchs/mux/-/badge.svg)](https://sourcegraph.com/github.com/gbrlsnchs/mux?badge)
[![GoDoc](https://godoc.org/github.com/gbrlsnchs/mux?status.svg)](https://godoc.org/github.com/gbrlsnchs/mux)
[![Minimal version](https://img.shields.io/badge/minimal%20version-go1.10%2B-5272b4.svg)](https://golang.org/doc/go1.10)

## About
This package is a fast HTTP multiplexer.

It uses a radix tree to match URLs. When matching simple routes, it's a zero allocation search.
It's fast, simple and supports middlewares in an elegant way.

## Usage
Full documentation [here].

### Installing
#### Go 1.10
`vgo get -u github.com/gbrlsnchs/mux`
#### Go 1.11 or after
`go get -u github.com/gbrlsnchs/mux`

### Importing
```go
import (
	// ...

	"github.com/gbrlsnchs/mux"
)
```

### Setting a handler (or handler function)
#### First, set a context key
```go
type key uint8

const ctxKey key = 0
```

#### Then, create a new router and set an endpoint handler
```go
rt := mux.NewRouter("/api", ctxKey)
rt.HandleFunc(http.MethodGet, "/ping", func(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
})
```

### Setting a common middleware for every endpoint
#### First, define a middleware
```go
func loggingFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
```

#### Then, create a router and use the middleware in all requests
```go
rt := mux.NewRouter("/api")
rt.Use(loggingFunc)
rt.Handle(http.MethodGet, "/ping", func(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
})
```

### Setting isolated middlewares
#### First, define a handler and some middlewares
```go
func handler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
})

func authFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if mux.Params(r.Context(), ctxKey)["secret"] != "my_secret" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func permissionFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !userIsAdmin(r) { // hypothetical function
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```

#### Then, create a middleware chain and add it to the router
```go
rt := mux.NewRouter("/", ctxKey)
guard := mux.NewChain(authFunc, permissionFunc)

rt.Handle(http.MethodPost, "/unprotected", handler)
rt.Handle(http.MethodPost, "/protected/:secret", guard(handler))
```

## Contributing
### How to help
- For bugs and opinions, please [open an issue](https://github.com/gbrlsnchs/mux/issues/new)
- For pushing changes, please [open a pull request](https://github.com/gbrlsnchs/mux/compare)