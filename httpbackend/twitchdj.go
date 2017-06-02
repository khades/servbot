package httpbackend

import (
	"encoding/json"
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func putTwitchDJ(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	decoder := json.NewDecoder(r.Body)
	var request models.TwitchDJ
	err := decoder.Decode(&request)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	twitchDJUpdateRequest := models.TwitchDJ{ID: request.ID, NotifyOnChange: request.NotifyOnChange}
	repos.PushTwitchDJ(channelID, &twitchDJUpdateRequest)
	json.NewEncoder(w).Encode(optionResponse{"OK"})
}
