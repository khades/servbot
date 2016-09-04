package repos

import "github.com/khades/servbot/models"

// GetChannelAliasedTemplates is setting
func GetChannelAliasedTemplates(channel string, aliasTo string) ([]models.TemplateInfo, error) {
	var result []models.TemplateInfo
	error := Db.C("templates").Find(models.TemplateAliasSelector{Channel: channel, AliasTo: aliasTo}).All(&result)
	return result, error
}
