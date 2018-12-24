package songRequestAPI

import (
	"encoding/json"
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/songRequest"
	"net/http"
	"strconv"

	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/eventbus"
	"github.com/khades/servbot/httpSession"

	"goji.io/pat"
)

type Service struct {
	songRequestService *songRequest.Service
	httpAPIService *httpAPI.Service
}

func (service *Service) get(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	value := service.songRequestService.Get(&channelInfo.ChannelID)

	//token, _ := service.songRequestServiceGetChannelToken(*&channelInfo.ChannelID)
	result := songRequestResponse{
		value, 
		channelInfo.GetIfUserIsMod(&s.UserID), 
		channelInfo.ChannelID == s.UserID, 
		//token,
	}
	json.NewEncoder(w).Encode(&result)
}

// func songrequestsWidget(w http.ResponseWriter, r *http.Request, &channelInfo.ChannelID *string) {
// 	value := service.songRequestServiceGetSongRequest(&channelInfo.ChannelID)
// 	json.NewEncoder(w).Encode(&value)

// }

func (service *Service) skip(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	id := pat.Param(r, "videoID")
	if id == "" {
		httpAPI.WriteJSONError(w, "song id is not defined", http.StatusNotFound)
		return
	}
	service.songRequestService.Pull(&channelInfo.ChannelID, &id)
}

func (service *Service) bubbleUp(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	id := pat.Param(r, "videoID")
	if id == "" {
		httpAPI.WriteJSONError(w, "song id is not defined", http.StatusNotFound)
		return
	}
	found := service.songRequestService.BubbleUp(&channelInfo.ChannelID, &id)
	json.NewEncoder(w).Encode(found)

}

func (service *Service) bubbleUpToSecond(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	id := pat.Param(r, "videoID")
	if id == "" {
		httpAPI.WriteJSONError(w, "song id is not defined", http.StatusNotFound)
		return
	}
	found := service.songRequestService.BubbleUpToSecond(&channelInfo.ChannelID, &id)
	json.NewEncoder(w).Encode(found)
}

func (service *Service) events(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	service.httpAPIService.WSEvent(w, r, eventbus.Songrequest(&channelInfo.ChannelID))
}

// func songrequestsWidgetEvents(w http.ResponseWriter, r *http.Request, &channelInfo.ChannelID *string) {
// 	websocketEventbusWriter(w, r, eventbus.Songrequest(&channelInfo.ChannelID))

// }
func (service *Service) setSettings(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	decoder := json.NewDecoder(r.Body)
	var request songRequest.ChannelSongRequestSettings
	err := decoder.Decode(&request)
	if err != nil {
		httpAPI.WriteJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	service.songRequestService.SetSettings(&channelInfo.ChannelID, &request)
}

func (service *Service) setVolume(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	volumeStr := pat.Param(r, "volume")
	volume, volumeError := strconv.Atoi(volumeStr)
	if volumeError != nil {
		httpAPI.WriteJSONError(w, volumeError.Error(), http.StatusUnprocessableEntity)
		return
	}
	if volume > 100 || volume < 0 {
		httpAPI.WriteJSONError(w, "Invalid value", http.StatusUnprocessableEntity)
		return
	}
	service.songRequestService.SetVolumeNoEvent(&channelInfo.ChannelID, volume)
}
