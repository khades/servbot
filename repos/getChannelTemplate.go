package repos

import "github.com/khades/servbot/models"

var templateCollection = "templates"

// GetChannelTemplate gets specific template for specific command on specific channel
func GetChannelTemplate(channelID *string, commandName *string) (models.TemplateInfo, error) {
	var result models.TemplateInfo
	error := Db.C(templateCollection).Find(models.TemplateSelector{ChannelID: *channelID, CommandName: *commandName}).One(&result)
	return result, error
}
