package repos

import (
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

// GetChannelActiveTemplates gets all non-empty templates for channel
func GetChannelActiveTemplates(channelID *string) (*[]models.TemplateInfo, error) {
	var result []models.TemplateInfo
	error := Db.C(templateCollection).Find(bson.M{"channelid": *channelID, "template": bson.M{"$ne": ""}}).All(&result)
	return &result, error
}
