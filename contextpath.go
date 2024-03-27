package entando

import (
	"net/http"
	"os"
	"strings"
)

// ContextPath is a middleware built to handle the special environment variable
// injected by the Entando deployer. The context path is used to access the MS.
//
// The handled env variable is `SERVER_SERVLET_CONTEXT_PATH`. When empty, `h` is
// left untouched.
func ContextPath(h http.Handler) http.Handler {
	pathPrefix := os.Getenv("SERVER_SERVLET_CONTEXT_PATH")
	if pathPrefix == "" {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newURL := *r.URL
		newURL.Path = strings.TrimPrefix(newURL.Path, pathPrefix)
		r.URL = &newURL
		h.ServeHTTP(w, r)
	})
}
