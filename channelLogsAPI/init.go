package channelLogsAPI

import (
	"github.com/khades/servbot/channelLogs"
	"github.com/khades/servbot/httpAPI"
	"goji.io/pat"
)


func Init(httpAPIService *httpAPI.Service,channelLogsService *channelLogs.Service) {
	service := Service{channelLogsService}
	mux := httpAPIService.NewMux()

	mux.HandleFunc(pat.Get("/api/channel/:channel/logs"), httpAPIService.WithMod(service.getUsers))
	mux.HandleFunc(pat.Get("/api/channel/:channel/logs/search/:search"), httpAPIService.WithMod(service.searchUsers))
	mux.HandleFunc(pat.Get("/api/channel/:channel/logs/username/:user"), httpAPIService.WithMod(service.getByUsername))
	mux.HandleFunc(pat.Get("/api/channel/:channel/logs/userid/:userID"), httpAPIService.WithMod(service.getByUserid))



}
