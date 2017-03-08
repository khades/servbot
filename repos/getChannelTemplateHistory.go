package repos

import "github.com/khades/servbot/models"

// GetChannelTemplateWithHistory gets specific paginated
func GetChannelTemplateWithHistory(channelID *string, commandName *string) (*models.TemplateInfoWithHistory, error) {
	var result models.TemplateInfoWithHistory
	error := Db.C(templateCollection).Find(models.TemplateSelector{ChannelID: *channelID, CommandName: *commandName}).One(&result)
	return &result, error
}
