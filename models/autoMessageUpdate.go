package models

type AutoMessageUpdate struct {
	ID            string
	User          string `valid:"required"`
	UserID        string `valid:"required"`
	ChannelID     string `valid:"required"`
	Game          string `json:"game"`
	Message       string `valid:"required" json:"message"`
	MessageLimit  int    `valid:"required" json:"messageLimit"`
	DurationLimit int    `valid:"required" json:"durationlimit"`
}
