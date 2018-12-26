package pubsub

import (
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/channelLogs"
	"github.com/khades/servbot/config"
	"sync"
)

func Init(channelInfoService *channelInfo.Service,
	config *config.Config,
	channelLogsService *channelLogs.Service,
	wg *sync.WaitGroup) *Client {

	client := &Client{
		IsWorking:          false,
		channelInfoService: channelInfoService,
		config:             config,
		channelLogsService: channelLogsService,
	}
	
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		for {
			client.run()
			wg.Done()
		}
	}(wg)

	return client
}
