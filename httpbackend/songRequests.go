package httpbackend

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/khades/servbot/eventbus"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
	"goji.io/pat"
)

type songRequest struct {
	*models.ChannelSongRequest
	IsMod   bool `json:"isMod"`
	IsOwner bool `json:"isOwner"`
	Token string `json:"token"`
}

func songrequests(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	value := repos.GetSongRequest(channelID)
	channelInfo, error := repos.GetChannelInfo(channelID)
	if error != nil {
		writeJSONError(w, "That channel is not defined", http.StatusForbidden)
		return
	}
	token, _ := repos.GetChannelToken(*channelID)
	result := songRequest{value, channelInfo.GetIfUserIsMod(&s.UserID), *channelID == s.UserID, token}
	json.NewEncoder(w).Encode(&result)
}
func songrequestsWidget(w http.ResponseWriter, r *http.Request, channelID *string) {
	value := repos.GetSongRequest(channelID)
	json.NewEncoder(w).Encode(&value)

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
	json.NewEncoder(w).Encode(found)

}

func songrequestsBubbleUpToSecond(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	id := pat.Param(r, "videoID")
	if id == "" {
		writeJSONError(w, "song id is not defined", http.StatusNotFound)
		return
	}
	found := repos.BubbleUpVideoToSecond(channelID, &id)
	json.NewEncoder(w).Encode(found)
}

func songrequestsEvents(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	websocketEventbusWriter(w, r, eventbus.Songrequest(channelID))
}

func songrequestsWidgetEvents(w http.ResponseWriter, r *http.Request, channelID *string) {
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

func songrequestSetVolume(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	volumeStr := pat.Param(r, "volume")
	volume, volumeError := strconv.Atoi(volumeStr)
	if volumeError != nil {
		writeJSONError(w, volumeError.Error(), http.StatusUnprocessableEntity)
		return
	}
	if volume > 100 || volume < 0 {
		writeJSONError(w, "Invalid value", http.StatusUnprocessableEntity)
		return
	}
	repos.SetSongRequestVolumeNoEvent(channelID, volume)
}
