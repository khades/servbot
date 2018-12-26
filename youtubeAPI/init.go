package youtubeAPIClient

import "github.com/khades/servbot/config"

func Init(config *config.Config) *YouTubeAPIClient {
	return &YouTubeAPIClient{config: config}
}
