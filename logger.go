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

	recorder := &statusRecorder{
		ResponseWriter: w,
		Status:         200,
	}

	if logger.Handler != nil {
		logger.Handler.ServeHTTP(recorder, r)
	}

	logger.Log.Printf(
		"%d %s %s %s",
		recorder.Status,
		r.RemoteAddr,
		r.Method,
		r.RequestURI,
	)
}

type statusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}
