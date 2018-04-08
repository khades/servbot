package httpbackend

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/khades/servbot/eventbus"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
	"goji.io/pat"
)

type songRequest struct {
	*models.ChannelSongRequest
	IsMod   bool `json:"isMod"`
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
}

func songrequestsBubbleUp(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	id := pat.Param(r, "videoID")
	if id == "" {
		writeJSONError(w, "song id is not defined", http.StatusNotFound)
		return
	}
	found := repos.BubbleUpVideo(channelID, &id)
	log.Println(found)
}

func songrequestsBubbleUpToSecond(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	id := pat.Param(r, "videoID")
	if id == "" {
		writeJSONError(w, "song id is not defined", http.StatusNotFound)
		return
	}
	found := repos.BubbleUpVideoToSecond(channelID, &id)
	log.Println(found)
}

func songrequestsEvents(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	websocketEventbusWriter(w, r, eventbus.Songrequest(channelID))
}

func songrequestsPushSettings(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) { 
	decoder := json.NewDecoder(r.Body)
	var request models.ChannelSongRequestSettings
	err := decoder.Decode(&request)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	repos.PushSongRequestSettings(channelID, &request)
}