package mux_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	. "github.com/gbrlsnchs/mux"
)

const paramsKey = "params"

var params Params

func TestMux(t *testing.T) {
	testCases := []struct {
		method   string
		requests map[string]int
		m        *Mux
		fns      []MiddlewareFunc
		handlers map[string]http.Handler
		params   map[string]Params
	}{
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
		{
			method: http.MethodGet,
			requests: map[string]int{
				"/": http.StatusOK,
			},
			m: New("/"),
			handlers: map[string]http.Handler{
				"/": http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
		},
		{
			method: http.MethodGet,
			requests: map[string]int{
				"/": http.StatusOK,
			},
			m: New(""),
			handlers: map[string]http.Handler{
				"/": http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
		},
		{
			method: http.MethodGet,
			requests: map[string]int{
				"/": http.StatusOK,
			},
			m: New("/"),
			handlers: map[string]http.Handler{
				"": http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
		},
		{
			method: http.MethodGet,
			requests: map[string]int{
				"/": http.StatusOK,
			},
			m: New(""),
			handlers: map[string]http.Handler{
				"": http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
		},
		{
			method: http.MethodGet,
			requests: map[string]int{
				"/123":          http.StatusOK,
				"/123/test/456": http.StatusNoContent,
			},
			m: New("/"),
			handlers: map[string]http.Handler{
				"/:test": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					params = r.Context().Value(paramsKey).(Params)
					w.WriteHeader(http.StatusOK)
				}),
				"/:test/test/:testing": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					params = r.Context().Value(paramsKey).(Params)
					w.WriteHeader(http.StatusNoContent)
				}),
			},
			params: map[string]Params{
				"/123":          Params{"test": "123"},
				"/123/test/456": Params{"test": "123", "testing": "456"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			params = nil
			tc.m.SetCtxKey(paramsKey)
			tc.m.Use(tc.fns...)

			for path, handler := range tc.handlers {
				tc.m.Handle(tc.method, path, handler)
			}
			for path, status := range tc.requests {
				t.Run(fmt.Sprintf("%s %s", tc.method, path), func(t *testing.T) {
					w := httptest.NewRecorder()
					r := httptest.NewRequest(tc.method, path, nil)

					tc.m.ServeHTTP(w, r)
					if want, got := status, w.Code; want != got {
						t.Errorf("want %d, got %d", want, got)
					}
					if want, got := len(tc.params[path]), len(params); want != got {

					}
					if len(tc.params) > 0 {
						if want, got := tc.params[path], params; !reflect.DeepEqual(want, got) {
							t.Errorf("want %#v, got %#v", want, got)
						}
					}
				})
			}
		})
	}
}
