package httpbackend

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type logsUserStruct struct {
	Channel string                `json:"channel"`
	Users   []models.ChannelUsers `json:"users"`
}

func logsUsers(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {

	channelUsers, error := repos.GetChannelUsers(channelID)
	if error != nil {
		log.Println(error)
		writeJSONError(w, error.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(logsUserStruct{*channelName, *channelUsers})
}
