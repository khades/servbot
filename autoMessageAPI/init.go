package autoMessageAPI

import (
	"github.com/khades/servbot/autoMessage"
	"github.com/khades/servbot/httpAPI"
	"goji.io/pat"
)

func Init(httpAPIService *httpAPI.Service, autoMessageService *autoMessage.Service) {
	service := &Service{autoMessageService}
	mux := httpAPIService.NewMux()

	mux.HandleFunc(pat.Get("/api/channel/:channel/automessages"), httpAPIService.WithMod(service.list))
	mux.HandleFunc(pat.Get("/api/channel/:channel/automessages/removeinactive"), httpAPIService.WithMod(service.removeInactive))
	mux.HandleFunc(pat.Get("/api/channel/:channel/automessages/:messageID"), httpAPIService.WithMod(service.get))
	mux.HandleFunc(pat.Post("/api/channel/:channel/automessages"), httpAPIService.WithMod(service.create))
	mux.HandleFunc(pat.Post("/api/channel/:channel/automessages/:id"), httpAPIService.WithMod(service.update))


}
