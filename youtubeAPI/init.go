package youtubeAPI

import "github.com/khades/servbot/config"

func Init(config *config.Config) *Client {
	return &Client{config: config}
}
