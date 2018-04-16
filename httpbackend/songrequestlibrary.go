package httpbackend

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/eventbus"
	"github.com/khades/servbot/l10n"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
	"goji.io/pat"
)

type videolibraryResponse struct {
	Count int                             `json:"count"`
	Items []models.SongRequestLibraryItem `json:"items"`
}

func songrequestGetLibrary(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {

	pageInt := 1
	log.Println(r.FormValue("page"))
	if r.FormValue("page") != "" {
		pageInt, _ = strconv.Atoi(r.FormValue("page"))
		if pageInt < 0 {
			pageInt = 1
		}
	}
	library, libraryLength, libraryError := repos.GetVideoLibraryItems(pageInt)
	if libraryError != nil {
		writeJSONError(w, libraryError.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(&videolibraryResponse{
		Count: libraryLength,
		Items: library})
}
func songrequestsUnban(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	videoID := pat.Param(r, "videoID")
	if videoID == "" {
		writeJSONError(w, "Invalid videoID", http.StatusUnprocessableEntity)
		return
	}
	repos.PullTagFromVideo(&videoID, *channelID+"-restricted")
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
		channelInfo, channelInfoError := repos.GetChannelInfo(&result.ChannelID)

		if channelInfoError != nil {
			continue
		}
		if result.RemovedYoutubeRestricted == true {
			eventbus.EventBus.Publish(eventbus.Songrequest(channelID), "youtuberestricted:"+result.Title)
			bot.IrcClientInstance.SendPublic(&models.OutgoingMessage{
				Channel: channelInfo.Channel,
				Body:    fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestPulledYoutubeRestricted, result.Title)})
		}
		if result.RemovedTwitchRestricted == true {
			eventbus.EventBus.Publish(eventbus.Songrequest(channelID), "twitchrestricted:"+result.Title)
			bot.IrcClientInstance.SendPublic(&models.OutgoingMessage{
				Channel: channelInfo.Channel,
				Body:    fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestPulledTwitchRestricted, result.Title)})
		}
		if result.RemovedChannelRestricted == true {
			eventbus.EventBus.Publish(eventbus.Songrequest(channelID), "channelrestricted:"+result.Title)
			bot.IrcClientInstance.SendPublic(&models.OutgoingMessage{
				Channel: channelInfo.Channel,
				Body:    fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestPulledChannelRestricted, result.Title)})
		}
		if result.RemovedTagRestricted == true {
			eventbus.EventBus.Publish(eventbus.Songrequest(channelID), "tagrestricted:"+result.Title)
			bot.IrcClientInstance.SendPublic(&models.OutgoingMessage{
				Channel: channelInfo.Channel,
				Body:    fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestPulledTagRestricted, result.Title, result.Tag)})
		}

	}

}

func songrequestGetBannedTracks(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {

	pageInt := 1
	log.Println(r.FormValue("page"))
	if r.FormValue("page") != "" {
		pageInt, _ = strconv.Atoi(r.FormValue("page"))
		if pageInt < 0 {
			pageInt = 1
		}
	}
	library, libraryLength, libraryError := repos.GetBannedTracksForChannel(channelID, pageInt)
	if libraryError != nil {
		writeJSONError(w, libraryError.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(&videolibraryResponse{
		Count: libraryLength,
		Items: library})
}
