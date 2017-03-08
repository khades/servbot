package repos

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/hoisie/mustache"
	"github.com/khades/servbot/models"
)

var subAlertCollection = "subAlert"

// SetSubAlert updates stream status (start of stream, topic of stream)
func SetSubAlert(user *string, userID *string, subAlert *models.SubAlert) error {
	ResubTemplateCache.Drop(&subAlert.ChannelID)
	template, error := mustache.ParseString(subAlert.ResubMessage)
	if error != nil {
		return error
	}

	ResubTemplateCache.put(&subAlert.ChannelID, template)
	Db.C(subAlertCollection).Upsert(models.ChannelSelector{ChannelID: subAlert.ChannelID}, bson.M{
		"$set": bson.M{
			"enabled":      subAlert.Enabled,
			"submessage":   subAlert.SubMessage,
			"resubmessage": subAlert.ResubMessage,
			"repeatbody":   subAlert.RepeatBody},
		"$push": bson.M{
			"history": bson.M{
				"$each": []models.SubAlertHistory{models.SubAlertHistory{
					User:         *user,
					UserID:       *userID,
					Date:         time.Now(),
					Enabled:      subAlert.Enabled,
					SubMessage:   subAlert.SubMessage,
					ResubMessage: subAlert.ResubMessage,
					RepeatBody:   subAlert.RepeatBody}},
				"$sort":  bson.M{"date": -1},
				"$slice": 10}}})

	return nil
}
