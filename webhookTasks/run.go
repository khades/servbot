package webhookTasks

import (
	"github.com/khades/servbot/webhook"
	"time"
)

func Run(webhookService *webhook.Service) *time.Ticker {
	period := time.Minute * 15
	ticker := time.NewTicker(period)
	webhookService.Subscribe(period)
	go func() {
		for {
			<-ticker.C
			webhookService.Subscribe(period)
		}
	}()
	return ticker
}

