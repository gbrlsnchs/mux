# mux (HTTP multiplexer for Go)
[![Build Status](https://travis-ci.org/gbrlsnchs/mux.svg?branch=master)](https://travis-ci.org/gbrlsnchs/mux)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/gbrlsnchs/mux)
[![Sourcegraph](https://sourcegraph.com/github.com/gbrlsnchs/mux/-/badge.svg)](https://sourcegraph.com/github.com/gbrlsnchs/mux?badge)
[![GoDoc](https://godoc.org/github.com/gbrlsnchs/mux?status.svg)](https://godoc.org/github.com/gbrlsnchs/mux)
[![Minimal version](https://img.shields.io/badge/minimal%20version-go1.10%2B-5272b4.svg)](https://golang.org/doc/go1.10)

## About
This is a complete rewrite of [httpmux], a now deprecated HTTP multiplexer.

This package uses a radix tree to match URLs. When it matches simple routes, it's a zero allocation search.
It's fast, simple and supports middlewares in an elegant way.

## Usage
Full documentation [here].

## Example
### Simple usage
```go
handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
})
m := mux.New("/api")
m.Handle(http.MethodGet, "/ping", handler)
```
### Common middleware
```go
loggingFunc := func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
m := mux.New("/api")
m.Use(loggingFunc)
m.Handle(http.MethodGet, "/ping", handler)
```
### Isolated middleware and URL parameters
```go
handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
})
authFunc := func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if params := r.Context().Value("my_params_key").(mux.Params); params.Get("secret") != "my_secret" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
guard := mux.Chain(authFunc, someOtherAuthFunc)
m := mux.New("/")
m.Handle(http.MethodGet, "/no-auth", handler)
m.Handle(http.MethodPost, "/needs-auth/:secret", guard(handler))
```

## Contribution
### How to help:
- Pull Requests
- Issues
- Opinions

[Go]: https://golang.org
[httpmux]: https://github.com/gbrlsnchs/httpmux
[here]: https://godoc.org/github.com/gbrlsnchs/mux