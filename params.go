package mux

// Params is a map of URL parameters.
type Params map[string]string

// Get retrieves a value and converts it to string.
func (path Params) Get(key string) string {
	return path[key]
}
