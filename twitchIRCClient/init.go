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
	wg *sync.WaitGroup) (*TwitchIRCClient, *time.Ticker, *time.Ticker) {

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

	ircClientTicker := time.NewTicker(time.Second * 3)

	go func(wg *sync.WaitGroup) {
		for {
			wg.Add(1)
			<-ircClientTicker.C
			twitchIRCClient.SendMessages(3)
			wg.Done()
		}
	}(wg)

	modTicker := time.NewTicker(time.Second * 10)

	go func(wg *sync.WaitGroup) {
		for {
			<-modTicker.C
			wg.Add(1)
			twitchIRCClient.SendModsCommand()
			wg.Done()
		}
	}(wg)

	return &twitchIRCClient, ircClientTicker, modTicker
}
