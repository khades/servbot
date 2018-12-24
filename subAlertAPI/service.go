package subAlertAPI

import (
	"encoding/json"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/httpSession"
	"github.com/khades/servbot/subAlert"
	"net/http"
)

type Service struct {
	subAlertService *subAlert.Service
}

func (service *Service) get(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	subBody, error := service.subAlertService.GetWithHistory(&channelInfo.ChannelID)
	if error != nil && error.Error() != "not found" {
		httpAPI.WriteJSONError(w, error.Error(), http.StatusUnprocessableEntity)
		return
	}
	json.NewEncoder(w).Encode(*subBody)
}

func (service *Service) set(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	decoder := json.NewDecoder(r.Body)
	var request subAlert.SubAlert
	err := decoder.Decode(&request)
	if err != nil {
		httpAPI.WriteJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	request.ChannelID = channelInfo.ChannelID
	validator := service.subAlertService.Set(&s.Username, &s.UserID, &request)
	if validator.Error == true {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(&httpAPI.HTTPError{Code: http.StatusUnprocessableEntity, Message: *validator})
		return
	}
	json.NewEncoder(w).Encode(httpAPI.OptionResponse{"OK"})
}
