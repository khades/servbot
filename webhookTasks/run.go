package webhookTasks

import (
	"time"

	"github.com/khades/servbot/webhook"
)

func Run(webhookService *webhook.Service) *time.Ticker {
	period := time.Minute * 15
	ticker := time.NewTicker(period)
	webhookService.Subscribe(period)
	go func() {
		for range ticker.C {
			webhookService.Subscribe(period)
		}
	}()
	return ticker
}
