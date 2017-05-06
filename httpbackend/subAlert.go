package httpbackend

import (
	"encoding/json"
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type subAlertResponse struct {
	Channel  string                     `json:"channel"`
	SubAlert models.SubAlertWithHistory `json:"subAlert"`
}

func subAlert(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	subBody, error := repos.GetSubAlertWithHistory(channelID)
	if error != nil && error.Error() != "not found" {
		writeJSONError(w, error.Error(), http.StatusUnprocessableEntity)
		return
	}
	json.NewEncoder(w).Encode(subAlertResponse{*channelName, *subBody})

}

func setSubAlert(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	decoder := json.NewDecoder(r.Body)
	var request models.SubAlert
	err := decoder.Decode(&request)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	request.ChannelID = *channelID
	validator := repos.SetSubAlert(&s.Username, &s.UserID, &request)
	if validator.Error == true {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(&models.HttpError{Code: http.StatusUnprocessableEntity, Message: *validator})
		return
	}
	json.NewEncoder(w).Encode(optionResponse{"OK"})
}
