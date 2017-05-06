package models

// SubAlertdescribes subscription alert
type SubAlert struct {
	ChannelID    string `json:"channelID"`
	Enabled      bool   `json:"enabled"`
	SubAlertBody `bson:",inline"`
}

type SubAlertBody struct {
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
