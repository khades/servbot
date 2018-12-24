package vkGroupAPI

import (
	"encoding/json"
	"net/http"

	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/httpSession"
)

type Service struct {
	channelInfoService *channelInfo.Service
}

func (service *Service) set(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfoStruct *channelInfo.ChannelInfo) {
	decoder := json.NewDecoder(r.Body)
	var request channelInfo.VkGroupInfo
	err := decoder.Decode(&request)
	if err != nil {
		httpAPI.WriteJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	vkUpdateRequest := channelInfo.VkGroupInfo{GroupName: request.GroupName, NotifyOnChange: request.NotifyOnChange}
	service.channelInfoService.PushVkGroupInfo(&channelInfoStruct.ChannelID, &vkUpdateRequest)
	json.NewEncoder(w).Encode(httpAPI.OptionResponse{"OK"})
}
