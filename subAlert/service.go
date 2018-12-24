package subAlert

import (
	"github.com/globalsign/mgo"
	"time"

	"github.com/cbroglie/mustache"
	"github.com/khades/servbot/models"
	"github.com/globalsign/mgo/bson"
)

type Service struct {
	collection *mgo.Collection
}

// Search returns subalert for specified channel
func (service *Service) Get(channelID *string) (*SubAlert, error) {
	var result SubAlert
	error := service.collection.Find(models.ChannelSelector{ChannelID: *channelID}).One(&result)
	return &result, error
}

// GetWithHistory returns subalert for specified channel with changelog
func (service *Service) GetWithHistory(channelID *string) (*SubAlertWithHistory, error) {
	var result SubAlertWithHistory
	error := service.collection.Find(models.ChannelSelector{ChannelID: *channelID}).One(&result)
	return &result, error
}


// Set updates stream status (start of stream, topic of stream)
func  (service *Service) Set(user *string, userID *string, subAlert *SubAlert) *SubAlertValidation {
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
		service.collection.Upsert(models.ChannelSelector{ChannelID: subAlert.ChannelID}, bson.M{
			"$set": &subAlert,
			"$push": bson.M{
				"history": bson.M{
					"$each": []SubAlertHistory{SubAlertHistory{
						*user,
						*userID,
						time.Now(),
						*subAlert}},
					"$sort":  bson.M{"date": -1},
					"$slice": 10}}})
	}
	return &templateValidation
}
