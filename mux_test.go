package mux_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/gbrlsnchs/mux"
	. "github.com/gbrlsnchs/mux/internal/mocks"
)

var params = make(Params)

func TestMux(t *testing.T) {
	testTable := []struct {
		method   string
		requests map[string]int
		m        *Mux
		fns      []MiddlewareFunc
		handlers map[string]http.Handler
		params   Params
	}{
		// #1: Simple match without mux prefix.
		{
			method: http.MethodGet,
			requests: map[string]int{
				"/test": http.StatusOK,
			},
			m: New("/"),
			handlers: map[string]http.Handler{
				"/test": http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
		},
		// #2: Simple match with mux prefix.
		{
			method: http.MethodGet,
			requests: map[string]int{
				"/test": http.StatusOK,
			},
			m: New("/test"),
			handlers: map[string]http.Handler{
				"": http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
		},
		// #3: Route with different method.
		{
			method: http.MethodPost,
			requests: map[string]int{
				"/test": http.StatusOK,
			},
			m: New("/test"),
			handlers: map[string]http.Handler{
				"": http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
		},
		// #4: Route not found without mux prefix.
		{
			method: http.MethodGet,
			requests: map[string]int{
				"/test": http.StatusNotFound,
			},
			m: New("/"),
			handlers: map[string]http.Handler{
				"/testing": http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
		},
		// #5: Route not found with mux prefix.
		{
			method: http.MethodGet,
			requests: map[string]int{
				"/test": http.StatusNotFound,
			},
			m: New("/testing"),
			handlers: map[string]http.Handler{
				"/testing": http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
		},
		// #6: Route with one middleware.
		{
			method: http.MethodGet,
			requests: map[string]int{
				"/test": http.StatusOK,
			},
			m: New("/"),
			fns: []MiddlewareFunc{
				func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						next.ServeHTTP(w, r)
					})
				},
			},
			handlers: map[string]http.Handler{
				"/test": http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
		},
		// #7: Route with two middlewares.
		{
			method: http.MethodGet,
			requests: map[string]int{
				"/test": http.StatusOK,
			},
			m: New("/"),
			fns: []MiddlewareFunc{
				func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						next.ServeHTTP(w, r)
					})
				},
				func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						next.ServeHTTP(w, r)
					})
				},
			},
			handlers: map[string]http.Handler{
				"/test": http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
		},
		// #8: Route with middleware termination.
		{
			method: http.MethodGet,
			requests: map[string]int{
				"/test": http.StatusBadRequest,
			},
			m: New("/"),
			fns: []MiddlewareFunc{
				func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						w.WriteHeader(http.StatusBadRequest)
					})
				},
			},
			handlers: map[string]http.Handler{
				"/test": http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
		},
	}

	for ttNum, tt := range testTable {
		tt.m.SetCtxKey(CtxKey)
		tt.m.Use(tt.fns...)

		for path, handler := range tt.handlers {
			tt.m.Handle(tt.method, path, handler)
		}
		for path, status := range tt.requests {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.method, path, nil)

			tt.m.ServeHTTP(w, r)
			if want, get := status, w.Code; want != get {
				t.Errorf("test #%d, handler \"%s\": want %d, got %d\n", ttNum+1, path, want, get)
			}
		}
	}
}
