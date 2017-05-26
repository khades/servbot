package repos

import (
	"time"

	"github.com/khades/servbot/models"

	"gopkg.in/mgo.v2/bson"
)

// PutChannelTemplate puts template into database
func putChannelTemplate(user *string, userID *string, channelID *string, commandName *string, template *models.TemplateInfoBody) {
	var templateHistory = models.TemplateHistory{
		TemplateInfoBody: *template,
		User:             *user,
		UserID:           *userID,
		Date:             time.Now(),
	}
	var arrayUpdateQuery = bson.M{
		"$set": template,
		"$push": bson.M{
			"history": bson.M{
				"$each":  []models.TemplateHistory{templateHistory},
				"$sort":  bson.M{"date": -1},
				"$slice": 10}}}

	Db.C(templateCollection).Upsert(
		models.TemplateSelector{ChannelID: *channelID, CommandName: *commandName},
		arrayUpdateQuery)
	if template.AliasTo == *commandName {
		Db.C(templateCollection).UpdateAll(
			bson.M{"channelid": *channelID, "aliasto": template.AliasTo, "commandname": bson.M{"$ne": *commandName}},
			arrayUpdateQuery)
	}
}
