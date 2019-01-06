package donationAPI

import "github.com/khades/servbot/donation"

type donationResult struct {
	Page      int `json:"page"`
	Donations []donation.Donation `json:"donations"`
}
