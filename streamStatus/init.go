package streamStatus

import (
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/config"
	"github.com/khades/servbot/gameResolve"
	"github.com/khades/servbot/twitchAPI"
)

func Init(config *config.Config,
	channelInfoService *channelInfo.Service,
	gameResolveService *gameResolve.Service,
	twitchAPIClient *twitchAPI.Client) *Service {
	service := Service{config: config, channelInfoService: channelInfoService, gameResolveService: gameResolveService, twitchAPIClient: twitchAPIClient}

	return &service
}
