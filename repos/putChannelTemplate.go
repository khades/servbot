package repos

import (
	"time"

	"github.com/khades/servbot/models"

	"gopkg.in/mgo.v2/bson"
)

// PutChannelTemplate puts template into database
func PutChannelTemplate(user *string, channel *string, commandName *string, aliasTo *string, template *string) {
	var templateInfo = models.TemplateInfo{
		CommandName: *commandName,
		Channel:     *channel,
		AliasTo:     *aliasTo,
		Template:    *template}

	Db.C("templates").Upsert(
		models.TemplateSelector{Channel: *channel, CommandName: *commandName},
		templateInfo)
	Db.C("templatesHistory").Insert(models.TemplateHistory{
		templateInfo,
		*user,
		time.Now()})
	if *aliasTo == *commandName {
		Db.C("templates").UpdateAll(
			models.TemplateAliasSelector{Channel: *channel, AliasTo: *aliasTo},
			bson.M{
				"$set": bson.M{
					"template": *template,
					"aliasto":  *aliasTo}})
	}
}
