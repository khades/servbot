package httpbackend

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"

	"goji.io/pat"
)

func logsUsers(w http.ResponseWriter, r *http.Request, s *models.HTTPSession) {
	channel := pat.Param(r, "channel")
	if channel == "" {
		http.Error(w, "Сhannel variable is not defined", http.StatusUnprocessableEntity)
		return
	}
	log.Println(channel)
	channelUsers, error := repos.GetChannelUsers(&channel)
	if error != nil {
		log.Println(error)
		http.Error(w, error.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(*channelUsers)
}
