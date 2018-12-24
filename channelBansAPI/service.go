package channelBansAPI

import (
	"encoding/json"
	"github.com/khades/servbot/channelBans"
	"net/http"

	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/httpSession"
)

type Service struct {
	*channelBans.Service
}

func (service *Service) get(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	results, error := service.Get(&channelInfo.ChannelID)
	if error != nil && error.Error() != "not found" {
		httpAPI.WriteJSONError(w, error.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(&results)
}
