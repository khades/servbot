package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)
// SubdayRecord describes vote of specific user
type SubdayRecord struct {
	User   string `json:"user"`
	UserID string `json:"userID"`
	Game   string `json:"game"`
}
// SubdayWinnersHistory describes history of winners picking on subday
type SubdayWinnersHistory struct {
	Date    time.Time      `json:"date"`
	Winners []SubdayRecord `json:"winners"`
}

//Subday fully describes one subday even on channel
type Subday struct {
	ID             bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	ChannelID      string                 `json:"channelID"`
	IsActive       bool                   `json:"isActive"`
	SubsOnly       bool                   `json:"subsOnly"`
	Name           string                 `json:"name"`
	Date           time.Time              `json:"date"`
	Votes          []SubdayRecord         `json:"votes"`
	Winners        []SubdayRecord         `json:"winners"`
	WinnersHistory []SubdayWinnersHistory `json:"winnersHistory"`
}

//SubdayNoWinners fully describes one subday even on channel WIHTOUT winners
type SubdayNoWinners struct {
	ID             bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	ChannelID      string                 `json:"channelID"`
	IsActive       bool                   `json:"isActive"`
	SubsOnly       bool                   `json:"subsOnly"`
	Name           string                 `json:"name"`
	Date           time.Time              `json:"date"`
	Votes          []SubdayRecord         `json:"votes"`
}

// SubdayList describes simplified version of subday, used when getting list of subdays
type SubdayList struct {
	ID             bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	ChannelID      string                 `json:"channelID"`
	IsActive       bool                   `json:"isActive"`
	Name           string                 `json:"name"`
	Date           time.Time              `json:"date"`
}