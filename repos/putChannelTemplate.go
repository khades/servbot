package repos

import (
	"time"

	"github.com/khades/servbot/models"

	"gopkg.in/mgo.v2/bson"
)

// PutChannelTemplate puts template into database
func PutChannelTemplate(user *string, userID *string, channelID *string, commandName *string, aliasTo *string, template *string) {
	var templateHistory = models.TemplateHistory{
		User:     *user,
		UserID:   *userID,
		Date:     time.Now(),
		AliasTo:  *aliasTo,
		Template: *template}
	var arrayUpdateQuery = bson.M{
		"$set": bson.M{
			"template": *template,
			"aliasto":  *aliasTo},
		"$push": bson.M{
			"history": bson.M{
				"$each":  []models.TemplateHistory{templateHistory},
				"$sort":  bson.M{"date": -1},
				"$slice": 5}}}

	Db.C(templateCollection).Upsert(
		models.TemplateSelector{ChannelID: *channelID, CommandName: *commandName},
		arrayUpdateQuery)
	if *aliasTo == *commandName {
		Db.C(templateCollection).UpdateAll(
			bson.M{"channelid": *channelID, "aliasto": *aliasTo, "commandname": bson.M{"$ne": *commandName}},
			arrayUpdateQuery)
	}
}
