package donationSource

import "time"

type DonationSources struct {
	ChannelID string `json:"channelID"`
	Yandex    DonationSource `json:"yandex"`
}

type DonationSource struct {
	Enabled        bool `json:"enabled"`
	Key            string 	`json:"-"`
	ExpirationDate time.Time `json:"expirationDate"`
	LastCheck      time.Time `json:"lastCheck"`
}
