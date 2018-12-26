package autoMessageAPI

import (
	"encoding/json"
	"net/http"

	"github.com/cbroglie/mustache"
	"github.com/globalsign/mgo/bson"
	"github.com/khades/servbot/autoMessage"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/httpSession"
	"github.com/sirupsen/logrus"
	"goji.io/pat"
)

type Service struct {
	autoMessageService *autoMessage.Service
}

func (service *Service) list(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	result, error := service.autoMessageService.List(&channelInfo.ChannelID)
	if error != nil && error.Error() != "not found" {
		httpAPI.WriteJSONError(w, error.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result)

}

func (service *Service) removeInactive(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	service.autoMessageService.RemoveInactive(&channelInfo.ChannelID)
}

func (service *Service) get(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	id := pat.Param(r, "messageID")

	if id == "" {
		httpAPI.WriteJSONError(w, "messageID or channel variable is not defined", http.StatusUnprocessableEntity)
		return
	}
	result, error := service.autoMessageService.Get(&id, &channelInfo.ChannelID)
	if error != nil {
		httpAPI.WriteJSONError(w, error.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(*result)
}

func (service *Service) create(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "httpbackend",
		"feature": "autoMessage",
		"action":  "create"})
	decoder := json.NewDecoder(r.Body)
	var request autoMessage.AutoMessageUpdate
	err := decoder.Decode(&request)
	if err != nil {
		logger.Info("Decoding Error: %s")
		httpAPI.WriteJSONError(w, "Invalid entry", http.StatusUnprocessableEntity)
		return
	}
	request.User = s.Username
	request.UserID = s.UserID

	request.ChannelID = *&channelInfo.ChannelID
	_, mustasheError := mustache.ParseString(request.Message)
	if mustasheError != nil {

		httpAPI.WriteJSONError(w, mustasheError.Error(), http.StatusUnprocessableEntity)
		return

	}
	id, validationError := service.autoMessageService.Create(&request)
	if validationError != nil {
		httpAPI.WriteJSONError(w, "Validation Failed", http.StatusUnprocessableEntity)
		return
	}
	json.NewEncoder(w).Encode(autoMessageCreationResponse{*id})

}

func (service *Service) update(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	id := pat.Param(r, "id")
	logger := logrus.WithFields(logrus.Fields{
		"package": "httpbackend",
		"feature": "autoMessage",
		"action":  "update"})
	if id == "" || bson.IsObjectIdHex(id) == false {
		httpAPI.WriteJSONError(w, "channel or id variable are not defined", http.StatusUnprocessableEntity)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var request autoMessage.AutoMessageUpdate
	err := decoder.Decode(&request)
	if err != nil {
		logger.Info("Decoding Error: %s")

		httpAPI.WriteJSONError(w, "Invalid entry", http.StatusUnprocessableEntity)
		return
	}
	_, mustasheError := mustache.ParseString(request.Message)
	if mustasheError != nil {

		httpAPI.WriteJSONError(w, mustasheError.Error(), http.StatusUnprocessableEntity)
		return

	}
	request.User = s.Username
	request.UserID = s.UserID
	request.ChannelID = *&channelInfo.ChannelID
	request.ID = id
	validationError := service.autoMessageService.Update(&request)
	if validationError != nil {
		httpAPI.WriteJSONError(w, "Validation Failed", http.StatusUnprocessableEntity)
		return
	}

}
