package vkGroupAPI

import (
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/httpAPI"
	"goji.io/pat"
)

func Init(httpAPIService *httpAPI.Service, channelInfoService *channelInfo.Service) {
	service := Service{channelInfoService}
	mux := httpAPIService.NewMux()
	mux.HandleFunc(pat.Post("/api/channel/:channel/externalservices/vk"), httpAPIService.WithMod(service.set))
	mux.HandleFunc(pat.Options("/api/channel/:channel/externalservices/vk"), httpAPIService.Options)
}
