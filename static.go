package forge

import (
	"net/http"
	"strings"
)

// Static servers static files without directory listings
type Static struct {
	FileSystem      http.FileSystem
	NotFoundHandler http.Handler
	fileServer      http.Handler
}

// ServerHTTP satisfies the http.Handler interface
func (static *Static) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if static.fileServer == nil {
		static.fileServer = http.FileServer(static.FileSystem)
	}

	requestedFileName := r.URL.Path

	requestingDirectory := strings.HasSuffix(requestedFileName, "/")
	if requestingDirectory {
		requestedFileName += "index.html"
	}

	if !static.fileExists(requestedFileName) {
		static.notFound(w, r)
		return
	}

	static.fileServer.ServeHTTP(w, r)
}

func (static *Static) notFound(w http.ResponseWriter, r *http.Request) {
	if static.NotFoundHandler != nil {
		static.NotFoundHandler.ServeHTTP(w, r)
		return
	}

	notFoundHander(w, r)
}

func (static *Static) fileExists(path string) bool {
	file, err := static.FileSystem.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	return true
}
