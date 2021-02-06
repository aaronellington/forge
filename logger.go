package forge

import (
	"log"
	"net/http"
	"os"
)

// Logger logs all request before passing off to the Handler
type Logger struct {
	Handler http.Handler
	Log     *log.Logger
}

// ServerHTTP satisfies the http.Handler interface
func (logger *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if logger.Log == nil {
		logger.Log = log.New(os.Stdout, "", log.LstdFlags)
	}

	logger.Log.Printf("%s", r.URL.Path)

	if logger.Handler != nil {
		logger.Handler.ServeHTTP(w, r)
	}
}
