package templateAPI

import (
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/template"
	"goji.io/pat"
)

var templateCollection = "list"

func Init(httpAPIService *httpAPI.Service, templateService *template.Service) {
	service := Service{templateService}
	mux := httpAPIService.NewMux()

	mux.HandleFunc(pat.Get("/api/channel/:channel/list"), httpAPIService.WithSessionAndChannel(service.list))

	mux.HandleFunc(pat.Get("/api/channel/:channel/list/:commandName"), httpAPIService.WithMod(service.get))
	mux.HandleFunc(pat.Post("/api/channel/:channel/list/:commandName"), httpAPIService.WithMod(service.set))
	mux.HandleFunc(pat.Options("/api/channel/:channel/list/:commandName"), httpAPIService.Options)

	mux.HandleFunc(pat.Post("/api/channel/:channel/list/:commandName/setAliasTo"), httpAPIService.WithMod(service.setAlias))
	mux.HandleFunc(pat.Options("/api/channel/:channel/list/:commandName/setAliasTo"), httpAPIService.Options)


}
