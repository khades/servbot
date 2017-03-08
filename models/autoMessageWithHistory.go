package models

type AutoMessageWithHistory struct {
	AutoMessage `bson:",inline"`
	History     []AutoMessageHistory `json:"history"`
}
