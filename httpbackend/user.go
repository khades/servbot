package httpbackend

import (
	"encoding/json"
	"net/http"

	"github.com/khades/servbot/models"
)

type userResponse struct {
	Username string
}

func user(w http.ResponseWriter, r *http.Request, s *models.HTTPSession) {
	json.NewEncoder(w).Encode(userResponse{s.Username})
}
