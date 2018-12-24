
package subAlertAPI

import (
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/subAlert"
	"goji.io/pat"
)

func Init(httpAPIService *httpAPI.Service, subAlertService *subAlert.Service) {
	service := Service{subAlertService}
	mux := httpAPIService.NewMux()

	mux.HandleFunc(pat.Get("/api/channel/:channel/subalert"), httpAPIService.WithMod(service.get))
	mux.HandleFunc(pat.Post("/api/channel/:channel/subalert"), httpAPIService.WithMod(service.set))
	mux.HandleFunc(pat.Options("/api/channel/:channel/subalert"), httpAPIService.Options)
}

