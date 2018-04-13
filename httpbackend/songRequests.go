package httpbackend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/l10n"

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
func songrequestSetTag(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	videoID := pat.Param(r, "videoID")
	if videoID == "" {
		writeJSONError(w, "Invalid videoID", http.StatusUnprocessableEntity)
		return
	}
	tag := pat.Param(r, "tag")
	if tag == "" {
		writeJSONError(w, "Invalid tag", http.StatusUnprocessableEntity)
		return
	}
	if tag == "youtuberestricted" && *channelID != s.UserID {
		writeJSONError(w, "That tag is restricted to channel owner", http.StatusUnprocessableEntity)
		return
	}

	results := repos.AddTagToVideo(&videoID, tag, s.UserID, strings.ToLower(s.Username))

	for _, result := range results {
		if result.RemovedTwitchRestricted == true || result.RemovedYoutubeRestricted == true {
			channelInfo, channelInfoError := repos.GetChannelInfo(&result.ChannelID)

			if channelInfoError != nil {
				continue
			}
			if result.RemovedYoutubeRestricted == true {
				bot.IrcClientInstance.SendPublic(&models.OutgoingMessage{
					Channel: channelInfo.Channel,
					Body:    fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestPulledYoutubeRestricted, result.Title)})
			}
			if result.RemovedTwitchRestricted == true {
				bot.IrcClientInstance.SendPublic(&models.OutgoingMessage{
					Channel: channelInfo.Channel,
					Body:    fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestPulledTwitchRestricted, result.Title)})
			}
			if result.RemovedChannelRestricted == true {
				bot.IrcClientInstance.SendPublic(&models.OutgoingMessage{
					Channel: channelInfo.Channel,
					Body:    fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestPulledChannelRestricted, result.Title)})
			}
			if result.RemovedTagRestricted == true {
				bot.IrcClientInstance.SendPublic(&models.OutgoingMessage{
					Channel: channelInfo.Channel,
					Body:    fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestPulledTagRestricted, result.Title, result.Tag)})
			}
		}

	}

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
