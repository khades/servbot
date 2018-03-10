package models

import "time"

// SubAlert describes subscription alert settings on channel
type SubAlert struct {
	ChannelID    string `json:"channelID"`
	Enabled      bool   `json:"enabled"`
	SubAlertBody `bson:",inline"`
}

// SubAlertBody struct descibes all template settings for subalert
type SubAlertBody struct {
	FollowerMessage        string `json:"followerMessage"`
	SubPrimeMessage        string `json:"subPrimeMessage"`
	SubFiveMessage         string `json:"subFiveMessage"`
	SubTenMessage          string `json:"subTenMessage"`
	SubTwentyFiveMessage   string `json:"subTwentyFiveMessage"`
	ResubPrimeMessage      string `json:"resubPrimeMessage"`
	ResubFiveMessage       string `json:"resubFiveMessage"`
	ResubTenMessage        string `json:"resubTenMessage"`
	ResubTwentyFiveMessage string `json:"resubTwentyFiveMessage"`
	ResubPrimeSmile        string `json:"resubPrimeSmile"`
	ResubFiveSmile         string `json:"resubFiveSmile"`
	ResubTenSmile          string `json:"resubTenSmile"`
	ResubTwentyFiveSmile   string `json:"resubTwentyFiveSmile"`
}

//SubAlertValidation struct describes validation errors for subalert templates
type SubAlertValidation struct {
	Error           bool `json:"error"`
	PrimeError      bool `json:"primeError"`
	FiveError       bool `json:"fiveError"`
	TenError        bool `json:"tenError"`
	TwentyFiveError bool `json:"twentyFiveError"`
}

// SubAlertHistory describes edit history of subalert
type SubAlertHistory struct {
	User     string    `json:"user"`
	UserID   string    `json:"userID"`
	Date     time.Time `json:"date"`
	SubAlert `bson:",inline"`
}

// SubAlertWithHistory subalert settings for a channel with edit history
type SubAlertWithHistory struct {
	SubAlert `bson:",inline"`
	History  []SubAlertHistory `json:"history"`
}
