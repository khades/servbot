package httpbackend

import (
	"encoding/json"
	"net/http"

	"goji.io/pat"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func logsUsers(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	search := ""

	channelUsers, error := repos.GetChannelUsers(channelID, &search)
	if error != nil {
		writeJSONError(w, error.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(channelUsers)
}

func logsUsersSearch(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	search := pat.Param(r, "search")

	channelUsers, error := repos.GetChannelUsers(channelID, &search)
	if error != nil {
		writeJSONError(w, error.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(channelUsers)
}
