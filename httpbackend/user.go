package httpbackend

import (
	"encoding/json"
	"net/http"
)

type userResponse struct {
	Username string
}

func user(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(userResponse{"khadesru"})
}
