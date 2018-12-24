package subdayAPI

import (
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/subday"
	"goji.io/pat"
)

func Init(httpAPIService *httpAPI.Service, subdayService *subday.Service) {
	service := Service{subdayService}
	mux := httpAPIService.NewMux()
	mux.HandleFunc(pat.Get("/api/channel/:channel/subdays"), httpAPIService.WithMod(service.list))

	mux.HandleFunc(pat.Post("/api/channel/:channel/subdays/new"), httpAPIService.WithMod(service.create))
	mux.HandleFunc(pat.Options("/api/channel/:channel/subdays/new"), httpAPIService.Options)

	mux.HandleFunc(pat.Get("/api/channel/:channel/subdays/:subdayID"), httpAPIService.WithSessionAndChannel(service.get))

	mux.HandleFunc(pat.Get("/api/channel/:channel/subdays/:subdayID/close"), httpAPIService.WithMod(service.close))
	mux.HandleFunc(pat.Options("/api/channel/:channel/subdays/:subdayID/close"), httpAPIService.Options)

	mux.HandleFunc(pat.Get("/api/channel/:channel/subdays/:subdayID/randomize"), httpAPIService.WithMod(service.randomize))
	mux.HandleFunc(pat.Options("/api/channel/:channel/subdays/:subdayID/randomize"), httpAPIService.Options)

	mux.HandleFunc(pat.Get("/api/channel/:channel/subdays/:subdayID/randomizeSubs"), httpAPIService.WithMod(service.randomizeSubs))
	mux.HandleFunc(pat.Options("/api/channel/:channel/subdays/:subdayID/randomizeSubs"), httpAPIService.Options)

	mux.HandleFunc(pat.Get("/api/channel/:channel/subdays/:subdayID/randomizeNonSubs"), httpAPIService.WithMod(service.randomizeNonSubs))
	mux.HandleFunc(pat.Options("/api/channel/:channel/subdays/:subdayID/randomizeNonSubs"), httpAPIService.Options)

	mux.HandleFunc(pat.Get("/api/channel/:channel/subdays/:subdayID/pullwinner/:user"), httpAPIService.WithMod(service.pullWinner))
	mux.HandleFunc(pat.Options("/api/channel/:channel/subdays/:subdayID/pullwinner/:user"), httpAPIService.Options)

}
