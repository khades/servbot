package configTasks

import (
	"time"

	"github.com/khades/servbot/config"
	"github.com/khades/servbot/twitchAPI"
)

func Run(twitchAPIClient *twitchAPI.Client,
	config *config.Config) {
	service := Service{twitchAPIClient,
		config}
	period := 10 * time.Minute
	service.UpdateAPIKey(period)
	ticker := time.NewTicker(period)
	go func() {
		<-ticker.C
		service.UpdateAPIKey(period)
	}()
}
