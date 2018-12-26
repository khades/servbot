package channelLogsAPI

import (
	"encoding/json"
	"github.com/khades/servbot/channelLogs"
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/httpSession"
	"net/http"

	"github.com/khades/servbot/channelInfo"
	"goji.io/pat"
)

type Service struct {
	channelLogsService *channelLogs.Service
}

func (service *Service) getUsers(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	search := ""

	channelUsers, error := service.channelLogsService.GetUsers(&channelInfo.ChannelID, &search)
	if error != nil {
		httpAPI.WriteJSONError(w, error.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(channelUsers)
}

func (service *Service) searchUsers(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	search := pat.Param(r, "search")

	channelUsers, error := service.channelLogsService.GetUsers(&channelInfo.ChannelID, &search)
	if error != nil {
		httpAPI.WriteJSONError(w, error.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(channelUsers)
}

func (service *Service) getByUsername(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	user := pat.Param(r, "user")
	if user == "" {
		httpAPI.WriteJSONError(w, "Bad request, no user set", http.StatusBadRequest)
		return
	}
	userLogs, error := service.channelLogsService.GetByUserName(&user, &channelInfo.ChannelID)
	if error != nil && error.Error() != "not found" {
		httpAPI.WriteJSONError(w, error.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(*userLogs)
}

func (service *Service) getByUserid(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	user := pat.Param(r, "userID")
	if user == "" {
		httpAPI.WriteJSONError(w, "Bad request, no user set", http.StatusBadRequest)
		return
	}
	userLogs, error := service.channelLogsService.GetByUserID(&user, &channelInfo.ChannelID)
	if error != nil && error.Error() != "not found" {
		httpAPI.WriteJSONError(w, error.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(*userLogs)
}
