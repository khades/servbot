package httpbackend

import (
	"encoding/json"
	"net/http"
	"time"
)

type timeResponse struct {
	Time time.Time `json:"time"`
}

func getTime(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(timeResponse{time.Now()})
}
