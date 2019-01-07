package donationSourceTasks

import (
	"time"

	"github.com/khades/servbot/donation"
	"github.com/khades/servbot/donationSource"
	"github.com/khades/servbot/yandexMoney"
	"github.com/sirupsen/logrus"
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
	logger := logrus.WithFields(logrus.Fields{
		"package": "donationSourceTasks",
		"action":  "processOneChannel"})
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
		if source.Yandex.LastCheck == transaction.DateTime {
			logger.Debug("Passing that transaction, its already processed")
			continue
		}
		service.donationService.SetPaid(transaction.Details, source.ChannelID)
	}
	if lastCheck.IsZero() {
		lastCheck = time.Now()
	}
	service.donationSourceService.UpdateYandexLastCheck(source.ChannelID, lastCheck)
}
