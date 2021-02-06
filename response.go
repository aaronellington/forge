package forge

import (
	"encoding/json"
	"net/http"
)

// Response is a basic response structure
type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ResponseText responds to an http.Request with a text body
func ResponseText(w http.ResponseWriter, statusCode int, body string) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(ResponseTextNotFound))
}

// ResponseJSON responds to an http.Request with a JSON body
func ResponseJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	encoder := json.NewEncoder(w)
	_ = encoder.Encode(v)
}
