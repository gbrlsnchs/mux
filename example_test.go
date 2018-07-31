package mux_test

import (
	"log"
	"net/http"

	"github.com/gbrlsnchs/mux"
)

func ExampleChain() {
	loggingFunc := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s\n", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
	authFunc := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if params, ok := r.Context().Value("params_ctx_key").(mux.Params); !ok || params.Get("secret") != "g0_v3g4n" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	chain := mux.Chain(loggingFunc, authFunc)
	m := mux.New("/api")
	m.Handle(http.MethodPost, "/top-secret/:password", chain(handler))
}

func ExampleMux() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	m := mux.New("/api")
	m.Handle(http.MethodGet, "/ping", handler)
}
