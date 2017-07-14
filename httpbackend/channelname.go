package httpbackend

import (
	"encoding/json"
	"net/http"

	"github.com/khades/servbot/models"
)

type channelNameStruct struct {
	Channel string `json:"channel"`
}

func channelName(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	json.NewEncoder(w).Encode(channelNameStruct{*channelName})
}
