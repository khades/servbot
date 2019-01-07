package paypalIPN

import (
	"github.com/khades/servbot/httpAPI"
	"goji.io/pat"
)

func Init(httpAPIService *httpAPI.Service) {
	mux := httpAPIService.NewMux()
	mux.HandleFunc(pat.Get("/paypal/ipn"), ipnGet)

	mux.HandleFunc(pat.Post("/paypal/ipn"), ipn)
}
