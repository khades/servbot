package repos

import (
	"time"

	"github.com/khades/servbot/models"

	"gopkg.in/mgo.v2/bson"
)

// PutChannelTemplate is setting
func PutChannelTemplate(user string, channel string, commandName string, template string) {
	Db.C("templates").Upsert(
		TemplateAmbiguousSelector(channel, commandName),
		bson.M{
			"$set": bson.M{
				"template": template},
			"$push": bson.M{
				"history": &models.TemplateHistoryItem{Template: template, User: user, Date: time.Now()}}})
}
