package subscriptionInfoAPI

import (
	"encoding/json"
	"github.com/asaskevich/EventBus"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/httpSession"
	"github.com/khades/servbot/subscriptionInfo"
	"net/http"
	"strconv"
	"time"

	"github.com/khades/servbot/eventbus"
)

type Service struct {
	subscriptionInfoService *subscriptionInfo.Service
	httpAPIService *httpAPI.Service
	eventBus EventBus.Bus
}


func (service *Service) get(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	dateLimit := r.URL.Query().Get("limit")
	if dateLimit == "" {
		result, _ := service.subscriptionInfoService.GetDefault(&channelInfo.ChannelID)
		json.NewEncoder(w).Encode(result)
		return
	}

	unixTime, error := strconv.ParseInt(dateLimit, 10, 64)
	if error != nil {
		httpAPI.WriteJSONError(w, error.Error(), http.StatusUnprocessableEntity)
		return
	}
	date := time.Unix(0, unixTime*int64(time.Millisecond))
	result, _ := service.subscriptionInfoService.Get(&channelInfo.ChannelID, date)
	json.NewEncoder(w).Encode(result)
}



func (service *Service) events(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	service.httpAPIService.WSEvent(w, r, eventbus.EventSub(&channelInfo.ChannelID))
}