package templateAPI

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/httpSession"
	"github.com/khades/servbot/template"
	"goji.io/pat"
)

type Service struct {
	*template.Service
}

func (service *Service) get(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {

	commandName := strings.ToLower(pat.Param(r, "commandName"))
	if commandName == "" {
		httpAPI.WriteJSONError(w, "URL is not valid", http.StatusBadRequest)
		return
	}

	result, _ := service.GetWithHistory(&channelInfo.ChannelID, &commandName)

	result.ChannelID = channelInfo.ChannelID
	result.CommandName = commandName

	json.NewEncoder(w).Encode(*result)
}

func (service *Service) set(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	decoder := json.NewDecoder(r.Body)
	var request template.TemplateInfoBody
	err := decoder.Decode(&request)
	if err != nil {
		httpAPI.WriteJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	commandName := strings.ToLower(strings.Join(strings.Fields(pat.Param(r, "commandName")), ""))
	if commandName == "" {
		httpAPI.WriteJSONError(w, "URL is not valid", http.StatusBadRequest)
		return
	}
	request.AliasTo = commandName

	templateError := service.Set(&s.Username, &s.UserID, &channelInfo.ChannelID, &commandName, &request)

	if templateError != nil {
		httpAPI.WriteJSONError(w, templateError.Error(), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w).Encode(httpAPI.OptionResponse{"OK"})

}

func (service *Service) setAlias(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	decoder := json.NewDecoder(r.Body)
	var request aliasToRequest
	err := decoder.Decode(&request)

	if err != nil {
		httpAPI.WriteJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	commandName := strings.ToLower(strings.Join(strings.Fields(pat.Param(r, "commandName")), ""))

	if commandName == "" {
		httpAPI.WriteJSONError(w, "URL is not valid", http.StatusBadRequest)
		return
	}
	lcaseAlias := strings.ToLower(strings.Join(strings.Fields(strings.ToLower(request.AliasTo)), ""))
	service.SetAlias(&s.Username, &s.UserID, &channelInfo.ChannelID, &commandName, &lcaseAlias)

	json.NewEncoder(w).Encode(httpAPI.OptionResponse{"OK"})

}

func (service *Service) list(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {

	templates, error := service.List(&channelInfo.ChannelID)
	if error != nil && error.Error() != "not found" {
		httpAPI.WriteJSONError(w, error.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(templatesResponse{templates, channelInfo.GetIfUserIsMod(&s.UserID)})
}
