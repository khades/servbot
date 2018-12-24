package subtrainAPI

import (
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/httpAPI"
	"goji.io/pat"
)

func Init(httpAPIService *httpAPI.Service, channelInfo *channelInfo.Service) {
	//mux.HandleFunc(pat.Get("/api/widget/subtrain"), withTokenSession(subtrainWidget))
	//mux.HandleFunc(pat.Get("/api/widget/subtrainEvents"), withTokenSession(subtrainWidgetEvents))

	service := Service{channelInfo}

	mux:= httpAPIService.NewMux()

	mux.HandleFunc(pat.Get("/api/channel/:channel/subtrain"), httpAPIService.WithMod(service.get))
	mux.HandleFunc(pat.Options("/api/channel/:channel/subtrain"), httpAPIService.Options)
	mux.HandleFunc(pat.Post("/api/channel/:channel/subtrain"), httpAPIService.WithMod(service.set))
	//mux.HandleFunc(pat.Get("/api/channel/:channel/subtrain/events"), withMod(subtrainEvents))

}