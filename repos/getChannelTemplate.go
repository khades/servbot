package repos

import "github.com/khades/servbot/models"

// GetChannelTemplate is setting
func GetChannelTemplate(channel string, commandName string) (models.TemplateInfo, error) {
	var result models.TemplateInfo
	error := Db.C("templates").Find(models.TemplateSelector{Channel: channel, CommandName: commandName}).One(&result)
	return result, error
}
