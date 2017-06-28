package models

// TwitchDJ describes information about twitchDJ service
type TwitchDJ struct {
	ID             string `json:"id"`
	Playing        bool   `json:"playing"`
	Track          string `json:"track"`
	NotifyOnChange bool   `json:"notifyOnChange"`
}
