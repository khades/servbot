package twitchIRCClient

import (
	"sync"
	"time"

	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/config"
)

// Start function dials up connection for chathandler
func Init(
	config *config.Config,
	channelInfoService *channelInfo.Service,
	handle TwitchIRCHandle,
	wg *sync.WaitGroup) *TwitchIRCClient {

	twitchIRCClient := TwitchIRCClient{
		config:             config,
		channelInfoService: channelInfoService,
		handle:             handle,
		ready:              false,
		modChannelIndex:    0,
		bounces:            make(map[string]time.Time),
		messagesSent:       0,
	}

	go func(wg *sync.WaitGroup) {
		wg.Add(1)
		twitchIRCClient.Start()
		wg.Done()
	}(wg)


	return &twitchIRCClient
}
