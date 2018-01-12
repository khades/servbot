package repos

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/cbroglie/mustache"
	"github.com/khades/servbot/models"
)

var subAlertCollection = "subAlert"

type SubAlertValidation struct {
	Error           bool `json:"error"`
	PrimeError      bool `json:"primeError"`
	FiveError       bool `json:"fiveError"`
	TenError        bool `json:"tenError"`
	TwentyFiveError bool `json:"twentyFiveError"`
}

// SetSubAlert updates stream status (start of stream, topic of stream)
func SetSubAlert(user *string, userID *string, subAlert *models.SubAlert) *SubAlertValidation {
	templateValidation := SubAlertValidation{false, false, false, false, false}

	_, primeError := mustache.ParseString(subAlert.ResubPrimeMessage)
	if primeError != nil {
		templateValidation.PrimeError = true
		templateValidation.Error = true
	}

	_, fiveError := mustache.ParseString(subAlert.ResubFiveMessage)
	if fiveError != nil {
		templateValidation.FiveError = true
		templateValidation.Error = true
	}

	_, tenError := mustache.ParseString(subAlert.ResubTenMessage)
	if tenError != nil {
		templateValidation.TenError = true
		templateValidation.Error = true
	}

	_, twentyFiveError := mustache.ParseString(subAlert.ResubTwentyFiveMessage)
	if twentyFiveError != nil {
		templateValidation.TwentyFiveError = true
		templateValidation.Error = true
	}
	if templateValidation.Error == false {
		Db.C(subAlertCollection).Upsert(models.ChannelSelector{ChannelID: subAlert.ChannelID}, bson.M{
			"$set": &subAlert,
			"$push": bson.M{
				"history": bson.M{
					"$each": []models.SubAlertHistory{models.SubAlertHistory{
						*user,
						*userID,
						time.Now(),
						*subAlert}},
					"$sort":  bson.M{"date": -1},
					"$slice": 10}}})
	}
	return &templateValidation
}
