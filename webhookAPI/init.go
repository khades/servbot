package webhookAPI

import (
	"github.com/khades/servbot/followers"
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/streamStatus"
	"github.com/khades/servbot/webhook"
	"goji.io/pat"
)

func Init(httpAPIService *httpAPI.Service,
	webHookService *webhook.Service,
streamStatusService *streamStatus.Service,
followersService *followers.Service) {

	service := Service{webHookService, streamStatusService, followersService}

	mux := httpAPIService.NewMux()

	mux.HandleFunc(pat.Get("/api/webhook/streams"), httpAPIService.CorsEnabled(service.verify))
	mux.HandleFunc(pat.Post("/api/webhook/streams"), httpAPIService.CorsEnabled(service.streamStatus))
	mux.HandleFunc(pat.Options("/api/webhook/streams"), httpAPIService.Options)

	mux.HandleFunc(pat.Get("/api/webhook/follows"), httpAPIService.CorsEnabled(service.verify))
	mux.HandleFunc(pat.Post("/api/webhook/follows"), httpAPIService.CorsEnabled(service.follows))
	mux.HandleFunc(pat.Options("/api/webhook/follows"), httpAPIService.Options)

}
