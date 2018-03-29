package httpbackend

import (
	"encoding/json"
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type channelInfoResponseStruct struct {
	Channel     string                 `json:"channel"`
	IsMod       bool                   `json:"isMod"`
	ModChannels []models.ChannelWithID `json:"modChannels"`
}

func channel(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {

	channelInfo, error := repos.GetChannelInfo(channelID)
	if error != nil {
		writeJSONError(w, error.Error(), http.StatusInternalServerError)
		return
	}

	modChannels, _ := repos.GetModChannels(&s.UserID)

	channelInfoResponse := channelInfoResponseStruct{
		Channel:     *channelName,
		IsMod:       channelInfo.GetIfUserIsMod(&s.UserID),
		ModChannels: modChannels}

	json.NewEncoder(w).Encode(&channelInfoResponse)
}
