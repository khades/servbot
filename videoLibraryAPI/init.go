package videoLibraryAPI

import (
	"github.com/asaskevich/EventBus"
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/songRequest"
	"github.com/khades/servbot/twitchIRC"
	"github.com/khades/servbot/videoLibrary"
	"goji.io/pat"
)

func Init(
	httpAPIService *httpAPI.Service,
	videoLibraryService *videoLibrary.Service,
	songRequestService *songRequest.Service,
	twitchIRCClient *twitchIRC.Client,
	eventBus EventBus.Bus) {
	service := Service{
		videoLibraryService,
		songRequestService,
		twitchIRCClient,
		eventBus,
	}
	mux:= httpAPIService.NewMux()

	mux.HandleFunc(pat.Get("/api/channel/:channel/songrequests/:videoID/settag/:tag"), httpAPIService.WithMod(service.setTag))
	mux.HandleFunc(pat.Options("/api/channel/:channel/songrequests/:videoID/settag/:tag"), httpAPIService.Options)

	mux.HandleFunc(pat.Get("/api/channel/:channel/songrequests/:videoID/unban"), httpAPIService.WithMod(service.unban))
	mux.HandleFunc(pat.Options("/api/channel/:channel/songrequests/:videoID/unban"),  httpAPIService.Options)

	mux.HandleFunc(pat.Get("/api/channel/:channel/songrequests/library/get"), httpAPIService.WithMod(service.get))
	mux.HandleFunc(pat.Get("/api/channel/:channel/songrequests/bannedtracks"), httpAPIService.WithMod(service.getBanned))

}
