package repos

import (
	"time"

	"github.com/khades/servbot/models"

	"gopkg.in/mgo.v2/bson"
)

// PutChannelTemplate is setting
func PutChannelTemplate(user string, channel string, commandName string, aliasTo string, template string) {
	Db.C("templates").Upsert(
		models.TemplateSelector{Channel: channel, CommandName: commandName},
		bson.M{
			"$set": bson.M{
				"template": template,
				"aliasto":  aliasTo},
			"$push": bson.M{
				"history": &models.TemplateHistoryItem{Template: template, User: user, Date: time.Now()}}})
}
