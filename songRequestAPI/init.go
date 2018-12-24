package songRequestAPI

import (
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/songRequest"
	"goji.io/pat"
)

func Init(	httpAPIService *httpAPI.Service,
	songRequestService *songRequest.Service) {
	service := Service{
		songRequestService,
		httpAPIService,
	}
	mux := httpAPIService.NewMux()
	mux.HandleFunc(pat.Get("/api/channel/:channel/songrequests"), httpAPIService.WithSessionAndChannel(service.get))

	mux.HandleFunc(pat.Options("/api/channel/:channel/songrequests"),  httpAPIService.Options)

	mux.HandleFunc(pat.Get("/api/channel/:channel/songrequests/skip/:videoID"), httpAPIService.WithMod(service.skip))
	mux.HandleFunc(pat.Get("/api/channel/:channel/songrequests/bubbleup/:videoID"), httpAPIService.WithMod(service.bubbleUp))
	mux.HandleFunc(pat.Get("/api/channel/:channel/songrequests/bubbleuptosecond/:videoID"), httpAPIService.WithMod(service.bubbleUpToSecond))

	mux.HandleFunc(pat.Get("/api/channel/:channel/songrequests/events"), httpAPIService.WithSessionAndChannel(service.events))
	mux.HandleFunc(pat.Post("/api/channel/:channel/songrequests/settings"), httpAPIService.WithMod(service.setSettings))
	mux.HandleFunc(pat.Options("/api/channel/:channel/songrequests/settings"), httpAPIService.Options)

	mux.HandleFunc(pat.Get("/api/channel/:channel/songrequests/setvolume/:volume"), httpAPIService.WithMod(service.setVolume))
	mux.HandleFunc(pat.Options("/api/channel/:channel/songrequests/setvolume/:volume"), httpAPIService.Options)



	//mux.HandleFunc(pat.Get("/api/widget/songrequests"), withTokenSession(songrequestsWidget))
	//mux.HandleFunc(pat.Get("/api/widget/songrequestsEvents"), withTokenSession(songrequestsWidgetEvents))
}