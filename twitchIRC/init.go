package twitchIRC

import (
	"github.com/khades/servbot/metrics"
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
	metrics *metrics.Service,
	wg *sync.WaitGroup) *Client {
	twitchIRCClient := Client{
		config:             config,
		channelInfoService: channelInfoService,
		handle:             handle,
		metrics:            metrics,
		ready:              false,
		modChannelIndex:    0,
		bounces:            make(map[string]time.Time),
		messagesSent:       0,
	}

	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		twitchIRCClient.Start()
		wg.Done()
	}(wg)

	return &twitchIRCClient
}
