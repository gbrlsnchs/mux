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

func TestRouter(t *testing.T) {
	testCases := []struct {
		method   string
		requests map[string]int
		rt       *Router
		fns      []MiddlewareFunc
		handlers map[string]http.Handler
		params   map[string]Params
	}{
		{
			method: http.MethodGet,
			requests: map[string]int{
				"/test": http.StatusOK,
			},
			rt: New("/"),
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
			rt: New("/test"),
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
			rt: New("/test"),
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
			rt: New("/"),
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
			rt: New("/testing"),
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
			rt: New("/"),
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
			rt: New("/"),
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
			rt: New("/"),
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
			rt: New("/"),
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
			rt: New(""),
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
			rt: New("/"),
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
			rt: New(""),
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
			rt: New("/"),
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
			tc.rt.SetDebug(true)
			tc.rt.SetCtxKey(paramsKey)
			tc.rt.Use(tc.fns...)

			for path, handler := range tc.handlers {
				tc.rt.Handle(tc.method, path, handler)
			}
			for path, status := range tc.requests {
				t.Run(fmt.Sprintf("%s %s", tc.method, path), func(t *testing.T) {
					w := httptest.NewRecorder()
					r := httptest.NewRequest(tc.method, path, nil)

					t.Log(tc.rt.String())
					tc.rt.ServeHTTP(w, r)
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
func TestSubrouter(t *testing.T) {
	testCases := []struct {
		method   string
		requests map[string]int
		rt       *Router
		path     string
		handlers map[string]http.Handler
	}{
		{
			method: http.MethodGet,
			requests: map[string]int{
				"/mux/router/test": http.StatusOK,
			},
			rt:   New("/mux"),
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
			rt:   New("/mux"),
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
			rt:   New("/"),
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
			rt := tc.rt.Router(tc.path)
			rt.SetDebug(true)
			for path, handler := range tc.handlers {
				rt.Handle(tc.method, path, handler)
			}
			for path, status := range tc.requests {
				w := httptest.NewRecorder()
				r := httptest.NewRequest(tc.method, path, nil)

				t.Log(tc.rt.String())
				tc.rt.ServeHTTP(w, r)
				if want, got := status, w.Code; want != got {
					t.Errorf("want %d, got %d", want, got)
				}
			}
		})
	}
}
