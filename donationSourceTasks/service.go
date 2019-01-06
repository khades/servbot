package donationSourceTasks

import (
	"time"

	"github.com/khades/servbot/donation"
	"github.com/khades/servbot/donationSource"
	"github.com/khades/servbot/yandexMoney"
)

type Service struct {
	donationService       *donation.Service
	donationSourceService *donationSource.Service
}

func (service *Service) Process() error {
	channels, error := service.donationSourceService.List()
	for _, channel := range channels {
		service.processOneChannel(channel)
	}
	return error
}

func (service *Service) processOneChannel(source donationSource.DonationSources) {
	if source.Yandex.Enabled == false {
		return
	}
	incomingTransactions, err := yandexMoney.GetHistory(source.Yandex.Key, source.Yandex.LastCheck)
	if err != nil {
		return
	}
	lastCheck := time.Time{}
	for _, transaction := range incomingTransactions.Operations {
		if transaction.DateTime.After(lastCheck) {
			lastCheck = transaction.DateTime
		}
		service.donationService.SetPaid(transaction.Details, source.ChannelID)
	}
	service.donationSourceService.UpdateYandexLastCheck(source.ChannelID, lastCheck)
}
