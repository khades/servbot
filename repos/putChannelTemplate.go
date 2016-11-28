package repos

import (
	"time"

	"github.com/khades/servbot/models"

	"gopkg.in/mgo.v2/bson"
)

// PutChannelTemplate puts template into database
func PutChannelTemplate(user *string, channel *string, commandName *string, aliasTo *string, template *string) {
	var templateHistory = models.TemplateHistory{
		User:     *user,
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
				"$slice": 10}}}

	Db.C("templates").Upsert(
		models.TemplateSelector{Channel: *channel, CommandName: *commandName},
		arrayUpdateQuery)
	if *aliasTo == *commandName {
		Db.C("templates").UpdateAll(
			models.TemplateAliasSelector{Channel: *channel, AliasTo: *aliasTo},
			arrayUpdateQuery)
	}
}
