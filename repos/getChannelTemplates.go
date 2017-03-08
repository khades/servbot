package repos

import "github.com/khades/servbot/models"

func GetChannelTemplates(channelID *string) (*[]models.TemplateInfo, error) {
	var result []models.TemplateInfo
	error := Db.C(templateCollection).Find(models.ChannelSelector{ChannelID: *channelID}).All(&result)
	return &result, error
}
