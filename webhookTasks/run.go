package webhookSchedule

import (
	"github.com/khades/servbot/webhook"
	"sync"
	"time"
)

func Init(webhookService *webhook.Service, wg *sync.WaitGroup) *time.Ticker {
	period := time.Minute * 15
	ticker := time.NewTicker(period)
	webhookService.Subscribe(period)
	go func(wg *sync.WaitGroup) {
		for {
			<-ticker.C
			wg.Add(1)
			webhookService.Subscribe(period)
			wg.Done()
		}
	}(wg)
	return ticker
}

