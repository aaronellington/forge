package forge

import "net/http"

// Security protects a http.Handler
type Security struct {
	EntryPoint EntryPoint
	Guards     []Guard
	Handler    http.Handler
}

// ServerHTTP satisfies the http.Handler interface
func (security *Security) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if security.Handler != nil {
		security.Handler.ServeHTTP(w, r)
	}
}

// EntryPoint defined the behavior when authentication is not present but required
type EntryPoint interface {
	Start(request http.Request)
}

// Guard that will attempt to authenticate a user
type Guard interface {
	Supports(request *http.Request) bool
	GetCredentials(request *http.Request) interface{}
	GetUser(credentials interface{}) (User, error)
	CheckCredentials(user interface{}, credentials interface{}) (bool, error)
	OnAuthenticationFailure(w http.ResponseWriter)
	OnAuthenticationSuccess(w http.ResponseWriter)
}

// User to be authenticated
type User interface {
	GetUsername() string
}
