package donationAPI

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/donation"
	"github.com/khades/servbot/httpSession"
)

type Service struct {
	donationService *donation.Service
}

func (service *Service) list(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	var page int = 1
	pageQuery := r.URL.Query().Get("page")
	if pageQuery == "" {
		pageResult, pageError := strconv.Atoi(pageQuery)
		if pageError != nil {
			page = pageResult
		}
	}
	donations, _ := service.donationService.List(channelInfo.ChannelID, page)
	result := &donationResult{Page: page, Donations: donations}
	json.NewEncoder(w).Encode(result)
}
