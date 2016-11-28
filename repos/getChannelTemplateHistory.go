package repos

import "github.com/khades/servbot/models"

// GetChannelTemplateWithHistory gets specific paginated
func GetChannelTemplateWithHistory(channel *string, commandName *string) (*models.TemplateInfoWithHistory, error) {
	var result models.TemplateInfoWithHistory
	error := Db.C("templates").Find(models.TemplateSelector{Channel: *channel, CommandName: *commandName}).One(&result)
	return &result, error
}
