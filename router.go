package forge

import (
	"net/http"
	"strings"
)

// Router serves http.Requests for a predefined map of Paths
type Router struct {
	NotFoundHander http.Handler
	routes         map[string]http.Handler
}

// ServerHTTP satisfies the http.Handler interface
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if stripTrailingSlash(w, r) {
		return
	}

	matchingRoute, found := router.routes[r.URL.Path]
	if !found {
		if router.NotFoundHander != nil {
			router.NotFoundHander.ServeHTTP(w, r)
			return
		}

		notFoundHander(w, r)
		return
	}

	if matchingRoute != nil {
		matchingRoute.ServeHTTP(w, r)
	}
}

// Handle registers a http.Handler to a predefined Path
func (router *Router) Handle(path string, handler http.Handler) {
	if router.routes == nil {
		router.routes = make(map[string]http.Handler)
	}

	router.routes[path] = handler
}

func notFoundHander(w http.ResponseWriter, r *http.Request) {
	RespondText(w, http.StatusNotFound, []byte(ResponseTextNotFound))
}

func stripTrailingSlash(w http.ResponseWriter, r *http.Request) bool {
	if !strings.HasSuffix(r.URL.Path, "/") || r.URL.Path == "/" {
		return false
	}

	http.Redirect(w, r, strings.TrimRight(r.URL.Path, "/"), http.StatusTemporaryRedirect)

	return true
}
