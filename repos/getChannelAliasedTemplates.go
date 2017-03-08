package repos

import "github.com/khades/servbot/models"

// GetChannelAliasedTemplates gets all templates that linked to a specified command, used when parent command gets updated
func GetChannelAliasedTemplates(channelID *string, aliasTo *string) ([]models.TemplateInfo, error) {
	var result []models.TemplateInfo
	error := Db.C(templateCollection).Find(models.TemplateAliasSelector{ChannelID: *channelID, AliasTo: *aliasTo}).All(&result)
	return result, error
}
