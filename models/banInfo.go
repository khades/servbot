package models

import "time"

// BanInfo struct describes info about person being banned
type BanInfo struct {
	User        string `json:"user"`
	Duration    int    `json:"duration"`
	Reason      string `json:"reason"`
	BanIssuer   string `json:"banIssuer"`
	BanIssuerID string `json:"banIssuerID"`
	Type string    `json:"type"`
	Date time.Time `json:"date"`
}
