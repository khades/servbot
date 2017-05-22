package models

import "time"

// BanInfo describes info about person being banned
type BanInfo struct {
	User     string    `json:"user"`
	Duration int       `json:"duration"`
	Reason   string    `json:"reason"`
	Type     string    `json:"type"`
	Date     time.Time `json:"date"`
}
