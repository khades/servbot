package httpbackend

import (
	"encoding/json"
	"net/http"

	"github.com/khades/servbot/eventbus"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
	"goji.io/pat"
)

type songRequest struct {
	*models.ChannelSongRequest
	IsMod bool `json:"isMod"`
	IsOwner bool `json:"isOwner"`
}

func songrequests(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	value := repos.GetSongRequest(channelID)
	channelInfo, error := repos.GetChannelInfo(channelID)
	if error != nil {
		writeJSONError(w, "That channel is not defined", http.StatusForbidden)
		return
	}

	result := songRequest{value, channelInfo.GetIfUserIsMod(&s.UserID), *channelID == s.UserID}
	json.NewEncoder(w).Encode(&result)
}

func songrequestsSkip(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	id := pat.Param(r, "videoID")
	if id == "" {
		writeJSONError(w, "song id is not defined", http.StatusNotFound)
		return
	}
	repos.PullSongRequest(channelID, &id)
	
	value := repos.GetSongRequest(channelID)
	
	channelInfo, error := repos.GetChannelInfo(channelID)
	if error != nil {
		writeJSONError(w, "That channel is not defined", http.StatusForbidden)
		return
	}

	result := songRequest{value, channelInfo.GetIfUserIsMod(&s.UserID), *channelID == s.UserID}
	json.NewEncoder(w).Encode(&result)
}

func songrequestsEvents(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	websocketEventbusWriter(w, r, eventbus.Songrequest(channelID))
}
