package models

// AutoMessageWithHistory struct descibes automessage with edit history
type AutoMessageWithHistory struct {
	AutoMessage `bson:",inline"`
	History     []AutoMessageHistory `json:"history"`
}
