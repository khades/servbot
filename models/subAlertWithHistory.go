package models

type SubAlertWithHistory struct {
	SubAlert `bson:",inline"`
	History  []SubAlertHistory `json:"history"`
}
