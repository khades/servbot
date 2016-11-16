package models

// SubAlertInfo describes subscription alert
type SubAlertInfo struct {
	Channel      string
	Enabled      bool
	SubMessage   string
	ResubMessage string
	RepeatBody   string
}
