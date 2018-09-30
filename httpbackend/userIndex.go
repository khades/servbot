package httpbackend

import (
	"encoding/json"
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type userIndexResponse struct {
	ModChannels []models.ChannelWithID `json:"modChannels"`
	Username    string                 `json:"username"`
	AvatarURL   string                 `json:"avatarUrl"`
	UserID      string                 `json:"userID"`
}

func userIndex(w http.ResponseWriter, r *http.Request, s *models.HTTPSession) {
	response := userIndexResponse{}
	channels, _ := repos.GetModChannels(&s.UserID)
	response.Username = s.Username
	response.AvatarURL = s.AvatarURL
	response.UserID = s.UserID
	response.ModChannels = channels
	json.NewEncoder(w).Encode(response)
}
