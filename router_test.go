package mux_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/gbrlsnchs/mux"
)

func TestRouter(t *testing.T) {
	testCases := []struct {
		method   string
		requests map[string]int
		m        *Mux
		path     string
		handlers map[string]http.Handler
	}{
		{
			method: http.MethodGet,
			requests: map[string]int{
				"/mux/router/test": http.StatusOK,
			},
			m:    New("/mux"),
			path: "/router",
			handlers: map[string]http.Handler{
				"/test": http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
		},
		{
			method: http.MethodPost,
			requests: map[string]int{
				"/mux/router/test": http.StatusOK,
			},
			m:    New("/mux"),
			path: "/router",
			handlers: map[string]http.Handler{
				"/test": http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
		},
		{
			method: http.MethodGet,
			requests: map[string]int{
				"/": http.StatusOK,
			},
			m:    New("/"),
			path: "/",
			handlers: map[string]http.Handler{
				"/": http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s %#v", tc.method, tc.requests), func(t *testing.T) {
			rt := tc.m.Router(tc.path)
			for path, handler := range tc.handlers {
				rt.Handle(tc.method, path, handler)
			}
			for path, status := range tc.requests {
				w := httptest.NewRecorder()
				r := httptest.NewRequest(tc.method, path, nil)

				tc.m.ServeHTTP(w, r)
				if want, got := status, w.Code; want != got {
					t.Errorf("want %d, got %d", want, got)
				}
			}
		})
	}
}
