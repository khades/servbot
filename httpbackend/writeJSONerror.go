package httpbackend

import (
	"encoding/json"
	"net/http"
)

func writeJSONError(w http.ResponseWriter, message string, headerCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(headerCode)
	json.NewEncoder(w).Encode(&httpError{Code: headerCode, Message: message})
}
