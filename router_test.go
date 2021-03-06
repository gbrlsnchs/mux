package mux_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	. "github.com/gbrlsnchs/mux"
)

type key uint8

const ctxKey key = 0

func TestRouter(t *testing.T) {
	var params map[string]string
	testCases := []struct {
		method      string
		requests    map[string]int
		rt          *Router
		fns         []MiddlewareFunc
		handlers    map[string]http.Handler
		params      map[string]map[string]string
		placeholder byte
	}{
		{
			method: http.MethodGet,
			requests: map[string]int{
				"/test": http.StatusOK,
			},
			rt: NewRouter("/", ctxKey),
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
			rt: NewRouter("/test", ctxKey),
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
			rt: NewRouter("/test", ctxKey),
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
			rt: NewRouter("/", ctxKey),
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
			rt: NewRouter("/testing", ctxKey),
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
			rt: NewRouter("/", ctxKey),
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
			rt: NewRouter("/", ctxKey),
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
			rt: NewRouter("/", ctxKey),
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
			rt: NewRouter("/", ctxKey),
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
			rt: NewRouter("", ctxKey),
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
			rt: NewRouter("/", ctxKey),
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
			rt: NewRouter("", ctxKey),
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
			rt: NewRouter("/", ctxKey),
			handlers: map[string]http.Handler{
				"/:test": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					params = Params(r.Context(), ctxKey)
					w.WriteHeader(http.StatusOK)
				}),
				"/:test/test/:testing": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					params = Params(r.Context(), ctxKey)
					w.WriteHeader(http.StatusNoContent)
				}),
			},
			params: map[string]map[string]string{
				"/123":          map[string]string{"test": "123"},
				"/123/test/456": map[string]string{"test": "123", "testing": "456"},
			},
		},
		{
			method: http.MethodGet,
			requests: map[string]int{
				"/api/user/123": http.StatusOK,
			},
			rt: NewRouter("/api", ctxKey),
			handlers: map[string]http.Handler{
				"/user/@id": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					params = Params(r.Context(), ctxKey)
					w.WriteHeader(http.StatusOK)
				}),
			},
			params: map[string]map[string]string{
				"/api/user/123": map[string]string{"id": "123"},
			},
			placeholder: '@',
		},
		{
			method: http.MethodGet,
			requests: map[string]int{
				"/api/user/@id": http.StatusOK,
			},
			rt: NewRouter("/api", ctxKey),
			handlers: map[string]http.Handler{
				"/user/@id": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					params = Params(r.Context(), ctxKey)
					w.WriteHeader(http.StatusOK)
				}),
			},
			params: map[string]map[string]string{
				"/api/user/@id": map[string]string{"id": "@id"},
			},
			placeholder: '@',
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			params = nil
			tc.rt.SetDebug(true)
			tc.rt.Use(tc.fns...)

			for path, handler := range tc.handlers {
				tc.rt.Handle(tc.method, path, handler)
			}
			for path, status := range tc.requests {
				t.Run(fmt.Sprintf("%s %s", tc.method, path), func(t *testing.T) {
					w := httptest.NewRecorder()
					r := httptest.NewRequest(tc.method, path, nil)

					t.Log(tc.rt.String())
					if tc.placeholder != 0 {
						tc.rt.SetPlaceholder(tc.placeholder)
					}
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
					t.Log(params)
				})
			}
		})
	}
}

func TestSubrouterUse(t *testing.T) {
	testCases := []struct {
		parentFn     MiddlewareFunc
		parentStatus int
		childFn      MiddlewareFunc
		childStatus  int
	}{
		{
			parentFn: func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					next.ServeHTTP(w, r)
				})
			},
			parentStatus: http.StatusOK,
			childFn: func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					next.ServeHTTP(w, r)
				})
			},
			childStatus: http.StatusOK,
		},
		{
			parentFn: func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
					next.ServeHTTP(w, r)
				})
			},
			parentStatus: http.StatusBadRequest,
			childFn: func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					next.ServeHTTP(w, r)
				})
			},
			childStatus: http.StatusBadRequest,
		},
		{
			parentFn: func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
					next.ServeHTTP(w, r)
				})
			},
			parentStatus: http.StatusBadRequest,
			childFn: func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadGateway)
					next.ServeHTTP(w, r)
				})
			},
			childStatus: http.StatusBadRequest,
		},
		{
			parentFn: func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					next.ServeHTTP(w, r)
				})
			},
			parentStatus: http.StatusOK,
			childFn: func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadGateway)
					next.ServeHTTP(w, r)
				})
			},
			childStatus: http.StatusBadGateway,
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			var (
				w *httptest.ResponseRecorder
				r *http.Request
			)
			parent := NewRouter("/parent", ctxKey)
			parent.Use(tc.parentFn)
			parent.HandleFunc(http.MethodGet, "", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("parent"))
			})

			child := parent.Router("/child")
			child.Use(tc.childFn)
			child.HandleFunc(http.MethodGet, "", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("child"))
			})

			w = httptest.NewRecorder()
			r = httptest.NewRequest(http.MethodGet, "/parent", nil)
			parent.ServeHTTP(w, r)
			if want, got := tc.parentStatus, w.Code; want != got {
				t.Errorf("want %d, got %d", want, got)
			}

			w = httptest.NewRecorder()
			r = httptest.NewRequest(http.MethodGet, "/parent/child", nil)
			child.ServeHTTP(w, r)
			if want, got := tc.childStatus, w.Code; want != got {
				t.Errorf("want %d, got %d", want, got)
			}
		})
	}
}

func TestSubrouter(t *testing.T) {
	var params map[string]string
	testCases := []struct {
		method      string
		requests    map[string]int
		rt          *Router
		path        string
		handlers    map[string]http.Handler
		params      map[string]map[string]string
		placeholder byte
	}{
		{
			method: http.MethodGet,
			requests: map[string]int{
				"/mux/router/test": http.StatusOK,
			},
			rt:   NewRouter("/mux", ctxKey),
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
			rt:   NewRouter("/mux", ctxKey),
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
			rt:   NewRouter("/", ctxKey),
			path: "/",
			handlers: map[string]http.Handler{
				"/": http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
		},
		{
			method: http.MethodGet,
			requests: map[string]int{
				"/api/user/123": http.StatusOK,
			},
			rt:   NewRouter("/api", ctxKey),
			path: "/user/@id",
			handlers: map[string]http.Handler{
				"": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					params = Params(r.Context(), ctxKey)
					w.WriteHeader(http.StatusOK)
				}),
			},
			params: map[string]map[string]string{
				"/api/user/123": map[string]string{"id": "123"},
			},
			placeholder: '@',
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s %#v", tc.method, tc.requests), func(t *testing.T) {
			if tc.placeholder > 0 {
				tc.rt.SetPlaceholder(tc.placeholder)
			}
			rt := tc.rt.Router(tc.path)
			rt.SetDebug(true)
			for path, handler := range tc.handlers {
				rt.Handle(tc.method, path, handler)
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
					if len(tc.params) > 0 {
						if want, got := tc.params[path], params; !reflect.DeepEqual(want, got) {
							t.Errorf("want %#v, got %#v", want, got)
						}
					}
					t.Log(params)
				})
			}
		})
	}
}
