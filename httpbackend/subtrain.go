package httpbackend

import (
	"encoding/json"
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func subtrain(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {

	channelInfo, error := repos.GetChannelInfo(channelID)
	if error != nil {
		writeJSONError(w, error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&channelInfo.SubTrain)
}

func putSubtrain(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	decoder := json.NewDecoder(r.Body)
	var request models.SubTrain
	err := decoder.Decode(&request)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	request.CurrentStreak = 0
	request.NotificationShown = false
	repos.PutChannelSubtrain(channelID, &request)
	json.NewEncoder(w).Encode(optionResponse{"OK"})
}
