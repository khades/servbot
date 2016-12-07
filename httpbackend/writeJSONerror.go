package httpbackend

import (
	"encoding/json"
	"net/http"

	"github.com/khades/servbot/models"
)

func writeJSONError(w http.ResponseWriter, message string, headerCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(headerCode)
	json.NewEncoder(w).Encode(&models.HttpError{Code: headerCode, Message: message})
}
