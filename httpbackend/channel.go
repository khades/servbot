package httpbackend

import (
	"encoding/json"
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type channelInfoResponseStruct struct {
	Channel string `json:"channel"`
	IsMod   bool   `json:"isMod"`
}

func channel(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {

	channelInfo, error := repos.GetChannelInfo(channelID)
	if error != nil {
		writeJSONError(w, error.Error(), http.StatusInternalServerError)
		return
	}
	channelInfoResponse := channelInfoResponseStruct{
		Channel: *channelName,
		IsMod:   channelInfo.GetIfUserIsMod(&s.UserID)}
	json.NewEncoder(w).Encode(&channelInfoResponse)
}
