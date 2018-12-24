package channelBansAPI

import (
	"github.com/khades/servbot/channelBans"
	"github.com/khades/servbot/httpAPI"
	"goji.io/pat"
)

func Init(httpAPIService *httpAPI.Service, channelBansService *channelBans.Service) {
	service := &Service{channelBansService}
	mux := httpAPIService.NewMux()
	mux.HandleFunc(pat.Get("/api/channel/:channel/bans"), httpAPIService.WithMod(service.get))
}
