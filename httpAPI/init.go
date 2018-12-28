package httpAPI

import (
	"github.com/asaskevich/EventBus"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/config"
	"github.com/khades/servbot/httpSession"
	"github.com/khades/servbot/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	goji "goji.io"
	"goji.io/pat"
)

func Init(config *config.Config,
	httpSessionService *httpSession.Service,
	channelInfoService *channelInfo.Service,
	eventBus EventBus.Bus,
	metrics *metrics.Service) *Service {
	service := &Service{
		config:             config,
		httpSessionService: httpSessionService,
		channelInfoService: channelInfoService,
		eventBus:           eventBus,
		requestsCounter:    make(map[string]requestCounterRecord),
		metrics:            metrics,
		mux:                goji.NewMux(),
	}

	service.mux.HandleFunc(pat.Get("/oauth"), service.oauth)
	service.mux.HandleFunc(pat.Get("/oauth/initiate"), service.oauthInitiate)
	service.mux.HandleFunc(pat.Get("/api/user"), service.WithAuth(service.userIndex))
	service.mux.HandleFunc(pat.Get("/api/channel/:channel"), service.WithSessionAndChannel(channelInfoHandler))
	service.mux.HandleFunc(pat.Get("/api/channel/:channel/channelname"), service.WithSessionAndChannel(channelName))
	service.mux.HandleFunc(pat.Get("/api/time"), service.CorsEnabled(getTime))
	service.mux.HandleFunc(pat.Get("/api/channel/:channel/info"), service.WithMod(channelInfoHandler))
	service.mux.Handle(pat.Get("/metrics"), promhttp.Handler())
	return service
}
