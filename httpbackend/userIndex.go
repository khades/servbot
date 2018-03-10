package httpbackend

import (
	"encoding/json"
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type userIndexResponse struct {
	ModChannels []models.ChannelWithID `json:"modChannels"`
}

func userIndex(w http.ResponseWriter, r *http.Request, s *models.HTTPSession) {
	response := userIndexResponse{}
	channels, _ := repos.GetModChannels(&s.UserID)
	response.ModChannels = channels
	json.NewEncoder(w).Encode(response)
}
