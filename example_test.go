package mux_test

import (
	"log"
	"net/http"

	"github.com/gbrlsnchs/mux"
)

func Example() {
	type key uint8
	const ctxKey key = 0

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	m := mux.NewRouter("/api", ctxKey)
	m.Handle(http.MethodGet, "/ping", handler)
}

func ExampleNewChain() {
	type key uint8
	const ctxKey key = 0

	m := mux.NewRouter("/api", ctxKey)
	guard := mux.NewChain(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Printf("%s %s", r.Method, r.URL.Path)
				next.ServeHTTP(w, r)
			})
		},
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if mux.Params(r.Context(), ctxKey)["secret"] != "my_secret" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				next.ServeHTTP(w, r)
			})
		},
	)

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	m.Handle(http.MethodPost, "/unprotected", handler)
	m.Handle(http.MethodPost, "/protected/:secret", guard(handler))
}
