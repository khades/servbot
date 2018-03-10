package repos

import (
	"time"

	"github.com/cbroglie/mustache"
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

var subAlertCollection = "subAlert"

// GetSubAlert returns subalert for specified channel
func GetSubAlert(channelID *string) (*models.SubAlert, error) {
	var result models.SubAlert
	error := db.C(subAlertCollection).Find(models.ChannelSelector{ChannelID: *channelID}).One(&result)
	return &result, error
}

// GetSubAlertWithHistory returns subalert for specified channel with changelog
func GetSubAlertWithHistory(channelID *string) (*models.SubAlertWithHistory, error) {
	var result models.SubAlertWithHistory
	error := db.C(subAlertCollection).Find(models.ChannelSelector{ChannelID: *channelID}).One(&result)
	return &result, error
}



// SetSubAlert updates stream status (start of stream, topic of stream)
func SetSubAlert(user *string, userID *string, subAlert *models.SubAlert) *models.SubAlertValidation {
	templateValidation := models.SubAlertValidation{false, false, false, false, false}

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
		db.C(subAlertCollection).Upsert(models.ChannelSelector{ChannelID: subAlert.ChannelID}, bson.M{
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
