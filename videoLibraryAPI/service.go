package videoLibraryAPI

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/asaskevich/EventBus"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/eventbus"
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/httpSession"
	"github.com/khades/servbot/l10n"
	"github.com/khades/servbot/songRequest"
	"github.com/khades/servbot/twitchIRC"
	"github.com/khades/servbot/videoLibrary"
	"goji.io/pat"
)

type Service struct {
	videoLibraryService *videoLibrary.Service
	songRequestService  *songRequest.Service
	twitchIRCClient     *twitchIRC.Client
	eventBus            EventBus.Bus
}

func (service *Service) get(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {

	pageInt := 1
	if r.FormValue("page") != "" {
		pageInt, _ = strconv.Atoi(r.FormValue("page"))
		if pageInt < 0 {
			pageInt = 1
		}
	}
	library, libraryLength, libraryError := service.videoLibraryService.List(pageInt)
	if libraryError != nil {
		httpAPI.WriteJSONError(w, libraryError.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(&videolibraryResponse{
		Count: libraryLength,
		Items: library})
}

func (service *Service) unban(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	videoID := pat.Param(r, "videoID")
	if videoID == "" {
		httpAPI.WriteJSONError(w, "Invalid videoID", http.StatusUnprocessableEntity)
		return
	}
	service.videoLibraryService.PullTag(&videoID, *&channelInfo.ChannelID+"-restricted")
}

func (service *Service) setTag(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	videoID := pat.Param(r, "videoID")
	if videoID == "" {
		httpAPI.WriteJSONError(w, "Invalid videoID", http.StatusUnprocessableEntity)
		return
	}
	tag := pat.Param(r, "tag")
	if tag == "" {
		httpAPI.WriteJSONError(w, "Invalid tag", http.StatusUnprocessableEntity)
		return
	}
	if tag == "youtuberestricted" && *&channelInfo.ChannelID != s.UserID {
		httpAPI.WriteJSONError(w, "That tag is restricted to channel owner", http.StatusUnprocessableEntity)
		return
	}

	results := service.songRequestService.PushTag(&videoID, tag, s.UserID, strings.ToLower(s.Username))

	for _, result := range results {

		if result.RemovedYoutubeRestricted == true {
			service.eventBus.Publish(eventbus.Songrequest(&channelInfo.ChannelID), "youtuberestricted:"+result.Title)
			service.twitchIRCClient.SendPublic(&twitchIRC.OutgoingMessage{
				Channel: channelInfo.Channel,
				Body:    fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestPulledYoutubeRestricted, result.Title)})
		}
		if result.RemovedTwitchRestricted == true {
			service.eventBus.Publish(eventbus.Songrequest(&channelInfo.ChannelID), "twitchrestricted:"+result.Title)
			service.twitchIRCClient.SendPublic(&twitchIRC.OutgoingMessage{
				Channel: channelInfo.Channel,
				Body:    fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestPulledTwitchRestricted, result.Title)})
		}
		if result.RemovedChannelRestricted == true {
			service.eventBus.Publish(eventbus.Songrequest(&channelInfo.ChannelID), "channelrestricted:"+result.Title)
			service.twitchIRCClient.SendPublic(&twitchIRC.OutgoingMessage{
				Channel: channelInfo.Channel,
				Body:    fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestPulledChannelRestricted, result.Title)})
		}
		if result.RemovedTagRestricted == true {
			service.eventBus.Publish(eventbus.Songrequest(&channelInfo.ChannelID), "tagrestricted:"+result.Title)
			service.twitchIRCClient.SendPublic(&twitchIRC.OutgoingMessage{
				Channel: channelInfo.Channel,
				Body:    fmt.Sprintf(l10n.GetL10n(channelInfo.Lang).SongRequestPulledTagRestricted, result.Title, result.Tag)})
		}

	}

}

func (service *Service) getBanned(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {

	pageInt := 1
	if r.FormValue("page") != "" {
		pageInt, _ = strconv.Atoi(r.FormValue("page"))
		if pageInt < 0 {
			pageInt = 1
		}
	}
	library, libraryLength, libraryError := service.videoLibraryService.ListBannedTracks(&channelInfo.ChannelID, pageInt)
	if libraryError != nil {
		httpAPI.WriteJSONError(w, libraryError.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(&videolibraryResponse{
		Count: libraryLength,
		Items: library})
}
