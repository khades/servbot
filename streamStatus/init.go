package streamStatus

import (
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/config"
	"github.com/khades/servbot/gameResolve"
	"github.com/khades/servbot/twitchAPIClient"
)

func Init(config *config.Config,
	channelInfoService *channelInfo.Service,
	gameResolveService *gameResolve.Service,
	twitchAPIClient *twitchAPIClient.TwitchAPIClient) *Service {
	service := Service{config: config, channelInfoService: channelInfoService, gameResolveService: gameResolveService, twitchAPIClient: twitchAPIClient}

	return &service
}
