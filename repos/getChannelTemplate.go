package repos

import "github.com/khades/servbot/models"

// GetChannelTemplate gets specific template for specific command on specific channel
func GetChannelTemplate(channel *string, commandName *string) (models.TemplateInfo, error) {
	var result models.TemplateInfo
	error := Db.C("templates").Find(models.TemplateSelector{Channel: *channel, CommandName: *commandName}).One(&result)
	return result, error
}
