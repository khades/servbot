package donationSourceAPI

import (
	"encoding/json"
	"net/http"

	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/donationSource"
	"github.com/khades/servbot/httpSession"
)

type Service struct {
	donationSourceService *donationSource.Service
}

func (service *Service) get(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	donationsSources, _ := service.donationSourceService.Get(channelInfo.ChannelID)
	json.NewEncoder(w).Encode(donationsSources)
}
