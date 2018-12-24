package subscriptionInfoAPI

import (
	"github.com/asaskevich/EventBus"
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/subscriptionInfo"
	"goji.io/pat"
)

func Init(
	httpAPIService *httpAPI.Service,
	subscriptionInfoService *subscriptionInfo.Service,
	eventBus EventBus.Bus) {
	service := Service{
		subscriptionInfoService,
		httpAPIService,
		eventBus,
	}

	mux := httpAPIService.NewMux()

	mux.HandleFunc(pat.Get("/api/channel/:channel/subs"), httpAPIService.WithMod(service.get))
	mux.HandleFunc(pat.Get("/api/channel/:channel/subs/events"), httpAPIService.WithMod(service.events))

}
