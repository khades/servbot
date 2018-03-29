package httpbackend

import (
	"encoding/json"
	"net/http"

	"github.com/khades/servbot/models"
)

type userResponse struct {
	Username  string `json:"username"`
	AvatarURL string `json:"avatarUrl"`
	UserID string `json:"userID"`
}

func user(w http.ResponseWriter, r *http.Request, s *models.HTTPSession) {
	json.NewEncoder(w).Encode(userResponse{s.Username, s.AvatarURL, s.UserID})
}
