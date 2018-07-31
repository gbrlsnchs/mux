/*
Package mux is an HTTP mux multiplexer that uses a radix tree as its backbone.

Create a main mux, which can contain a prefix:

	m := mux.New("/api")

Don't forget to set a context key (preferably a const iota).
It's used for accessing URL parameters:

	m.SetCtxKey("params_ctx_key")

If desired, set middlewares that'll run on every request:

	loggingFunc := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s\n", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
	m.Use(loggingFunc)

Then set a handler:

	m.Handle(http.MethodGet, "/ping", myHTTPHandler) // will be handled at "/api/ping"

If you want to use specific middlewares in certains endpoints,
simply create a middleware chain to pass as the handler:

	// Check credentials.
	authFunc := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			params, ok := r.Context().Value("params_ctx_key").(Params); !ok || params.Get("secret") != "g0_v3g4n" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
	// Check if user has needed permissions.
	permissionFunc := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if userKey, ok := r.Context().Value("user_ctx_key").(string); !ok || !isAdmin(userKey) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
	guard := mux.Chain(authFunc, permissionFunc)
	m.Handle(http.MethodPost, "/admin/:secret", guard(myHTTPHandler)) // will be handled at "/api/admin/:secret"
*/
package mux
