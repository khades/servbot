package template

import (
	"strings"
	"time"

	"github.com/cbroglie/mustache"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/khades/servbot/channelInfo"
)

type Service struct {
	collection         *mgo.Collection
	channelInfoService *channelInfo.Service
}

// SetAlias sets alias for a command by copying its template body and setting reference to another command for next updates
func (service *Service) SetAlias(user *string, userID *string, channelID *string, commandName *string, aliasTo *string) {
	commandNameFixed := strings.ToLower(strings.Join(strings.Fields(*commandName), ""))
	aliasFixed := strings.ToLower(strings.Join(strings.Fields(*aliasTo), ""))

	result, error := service.Get(channelID, &aliasFixed)
	aliasTemplate := TemplateInfoBody{}
	if error == nil {
		aliasTemplate = result.TemplateInfoBody
	}

	service.put(user, userID, channelID, &commandNameFixed, &aliasTemplate)
	service.pushCommandsForChannel(channelID)

}

// Set sets command template on channel
func (service *Service) Set(user *string, userID *string, channelID *string, commandName *string, template *TemplateInfoBody) error {
	commandNameFixed := strings.ToLower(strings.Join(strings.Fields(*commandName), ""))

	if template.Template == "" {
		_, templateError := mustache.ParseString(template.Template)
		if templateError != nil {
			return templateError
		}
	}
	service.put(user, userID, channelID, &commandNameFixed, template)
	service.pushCommandsForChannel(channelID)
	return nil
}

// Get returns template on specified channel for specified channel
func (service *Service) Get(channelID *string, commandName *string) (TemplateInfo, error) {
	var result TemplateInfo
	error := service.collection.Find(TemplateSelector{ChannelID: *channelID, CommandName: *commandName}).One(&result)
	return result, error
}

// GetWithHistory returns template and edit history on specified channel for specified channel
func (service *Service) GetWithHistory(channelID *string, commandName *string) (*TemplateInfoWithHistory, error) {
	var result TemplateInfoWithHistory
	error := service.collection.Find(TemplateSelector{ChannelID: *channelID, CommandName: *commandName}).One(&result)
	return &result, error
}

// List returns all temlates for specified channel
func (service *Service) List(channelID *string) ([]TemplateInfo, error) {
	var result []TemplateInfo
	error := service.collection.Find(bson.M{"channelid": *channelID}).All(&result)
	return result, error
}

// ListActive returns all commands with non-empty templates for specified channel
func (service *Service) ListActive(channelID *string) ([]TemplateInfo, error) {
	var result []TemplateInfo
	error := service.collection.Find(bson.M{"channelid": *channelID, "template": bson.M{"$ne": ""}}).All(&result)
	return result, error
}

// ListAliases returns all aliased commands for specified channel
func (service *Service) ListAliases(channelID *string, aliasTo *string) ([]TemplateInfo, error) {
	var result []TemplateInfo
	error := service.collection.Find(TemplateAliasSelector{ChannelID: *channelID, AliasTo: *aliasTo}).All(&result)
	return result, error
}

func (service *Service) pushCommandsForChannel(channelID *string) {
	activeTemplates, error := service.ListActive(channelID)
	if (error != nil) {

	}
	commands := []string{}
	for _, command := range activeTemplates {
		commands = append(commands, command.CommandName)
	}
	service.channelInfoService.SetCommandsForChannel(channelID, commands)
}

// PutChannelTemplate puts template into database, also storing edit history
func (service *Service) put(user *string, userID *string, channelID *string, commandName *string, template *TemplateInfoBody) {
	var templateHistory = TemplateHistory{
		TemplateInfoBody: *template,
		User:             *user,
		UserID:           *userID,
		Date:             time.Now(),
	}
	var arrayUpdateQuery = bson.M{
		"$set": template,
		"$push": bson.M{
			"history": bson.M{
				"$each":  []TemplateHistory{templateHistory},
				"$sort":  bson.M{"date": -1},
				"$slice": 10}}}

	service.collection.Upsert(
		TemplateSelector{ChannelID: *channelID, CommandName: *commandName},
		arrayUpdateQuery)
	if template.AliasTo == *commandName {
		service.collection.UpdateAll(
			bson.M{"channelid": *channelID, "aliasto": template.AliasTo, "commandname": bson.M{"$ne": *commandName}},
			arrayUpdateQuery)
	}
}
