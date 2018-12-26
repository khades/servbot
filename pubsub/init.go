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
		wg.Add(1)
	client := &Client{
		IsWorking:          false,
		channelInfoService: channelInfoService,
		config:             config,
		channelLogsService: channelLogsService,
	}

	go func(wg *sync.WaitGroup) {
		for {
			client.run()
			wg.Done()
		}
	}(wg)

	return client
}
