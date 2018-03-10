package repos

import (
	"time"

	"github.com/cbroglie/mustache"
	"github.com/khades/servbot/models"
	"gopkg.in/mgo.v2/bson"
)

var templateCollection = "templates"

// SetChannelTemplateAlias sets alias for a command by copying its template body and setting reference to another command for next updates
func SetChannelTemplateAlias(user *string, userID *string, channelID *string, commandName *string, aliasTo *string) {
	result, error := GetChannelTemplate(channelID, aliasTo)
	aliasTemplate := models.TemplateInfoBody{}
	if error == nil {
		aliasTemplate = result.TemplateInfoBody
	}

	putChannelTemplate(user, userID, channelID, commandName, &aliasTemplate)
	PushCommandsForChannel(channelID)

}

// SetChannelTemplate sets command template on channel
func SetChannelTemplate(user *string, userID *string, channelID *string, commandName *string, template *models.TemplateInfoBody) error {
	if template.Template == "" {
		_, templateError := mustache.ParseString(template.Template)
		if templateError != nil {
			return templateError
		}
	}
	putChannelTemplate(user, userID, channelID, commandName, template)
	PushCommandsForChannel(channelID)
	return nil
}

// GetChannelTemplate returns template on specified channel for specified channel
func GetChannelTemplate(channelID *string, commandName *string) (models.TemplateInfo, error) {
	var result models.TemplateInfo
	error := db.C(templateCollection).Find(models.TemplateSelector{ChannelID: *channelID, CommandName: *commandName}).One(&result)
	return result, error
}

// GetChannelTemplateWithHistory returns template and edit history on specified channel for specified channel
func GetChannelTemplateWithHistory(channelID *string, commandName *string) (*models.TemplateInfoWithHistory, error) {
	var result models.TemplateInfoWithHistory
	error := db.C(templateCollection).Find(models.TemplateSelector{ChannelID: *channelID, CommandName: *commandName}).One(&result)
	return &result, error
}

// GetChannelTemplates returns all temlates for specified channel
func GetChannelTemplates(channelID *string) ([]models.TemplateInfo, error) {
	var result []models.TemplateInfo
	error := db.C(templateCollection).Find(models.ChannelSelector{ChannelID: *channelID}).All(&result)
	return result, error
}

// GetChannelActiveTemplates returns all commands with non-empty templates for specified channel
func GetChannelActiveTemplates(channelID *string) ([]models.TemplateInfo, error) {
	var result []models.TemplateInfo
	error := db.C(templateCollection).Find(bson.M{"channelid": *channelID, "template": bson.M{"$ne": ""}}).All(&result)
	return result, error
}

// GetChannelAliasedTemplates returns all aliased commands for specified channel
func GetChannelAliasedTemplates(channelID *string, aliasTo *string) ([]models.TemplateInfo, error) {
	var result []models.TemplateInfo
	error := db.C(templateCollection).Find(models.TemplateAliasSelector{ChannelID: *channelID, AliasTo: *aliasTo}).All(&result)
	return result, error
}

// PutChannelTemplate puts template into database, also storing edit history
func putChannelTemplate(user *string, userID *string, channelID *string, commandName *string, template *models.TemplateInfoBody) {
	var templateHistory = models.TemplateHistory{
		TemplateInfoBody: *template,
		User:             *user,
		UserID:           *userID,
		Date:             time.Now(),
	}
	var arrayUpdateQuery = bson.M{
		"$set": template,
		"$push": bson.M{
			"history": bson.M{
				"$each":  []models.TemplateHistory{templateHistory},
				"$sort":  bson.M{"date": -1},
				"$slice": 10}}}

	db.C(templateCollection).Upsert(
		models.TemplateSelector{ChannelID: *channelID, CommandName: *commandName},
		arrayUpdateQuery)
	if template.AliasTo == *commandName {
		db.C(templateCollection).UpdateAll(
			bson.M{"channelid": *channelID, "aliasto": template.AliasTo, "commandname": bson.M{"$ne": *commandName}},
			arrayUpdateQuery)
	}
}
