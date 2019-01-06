package donationAPI

import (
	"github.com/khades/servbot/donation"
	"github.com/khades/servbot/httpAPI"
	"goji.io/pat"
)

func Init(httpAPIService *httpAPI.Service, donationService *donation.Service) {
	service := Service{donationService}
	mux := httpAPIService.NewMux()

	mux.HandleFunc(pat.Get("/api/channel/:channel/donations"), httpAPIService.WithOwner(service.list))
}
