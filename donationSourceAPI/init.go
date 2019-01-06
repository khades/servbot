package donationSourceAPI

import (
	"github.com/khades/servbot/donationSource"
	"github.com/khades/servbot/httpAPI"
	"goji.io/pat"
)

func Init(httpAPIService *httpAPI.Service, donationSourceService *donationSource.Service) {
	service := Service{donationSourceService}
	mux := httpAPIService.NewMux()

	mux.HandleFunc(pat.Get("/api/channel/:channel/donationsources"), httpAPIService.WithOwner(service.get))
}
