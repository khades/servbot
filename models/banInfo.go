package models

// BanInfo describes info about person being banned
type BanInfo struct {
	Duration int
	Reason   string `bson:",omitempty"`
}
