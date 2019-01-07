package donationSourceTasks

import (
	"time"

	"github.com/khades/servbot/donation"
	"github.com/khades/servbot/donationSource"
)

func Run(donationService *donation.Service,
	donationSourceService *donationSource.Service) {
	service := Service{donationService, donationSourceService}
	ticker := time.NewTicker(time.Second * 20)

	go func() {
		for range ticker.C {
			service.Process()
		}
	}()

}
