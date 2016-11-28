package repos

import "github.com/khades/servbot/models"

func GetChannelTemplates(channel *string) (*[]models.TemplateInfo, error) {
	var result []models.TemplateInfo
	error := Db.C("templates").Find(models.ChannelSelector{Channel: *channel}).All(&result)
	return &result, error
}
