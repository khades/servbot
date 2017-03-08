package models

// SubAlertdescribes subscription alert
type SubAlert struct {
	ChannelID    string `json:"channelID"`
	Enabled      bool   `json:"enabled"`
	SubMessage   string `json:"subMessage"`
	ResubMessage string `json:"resubMessage"`
	RepeatBody   string `json:"repeatBody"`
}
