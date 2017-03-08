package repos

import (
	"log"

	"github.com/cbroglie/mustache"
)

type templateContainer struct {
	templateMap map[string]*mustache.Template
	aliases     map[string]string
}

func (template templateContainer) Get(channelID *string, inputCommandName *string) (*mustache.Template, bool) {
	commandName := *inputCommandName
	alias, exists := template.aliases[*inputCommandName]
	if exists {
		commandName = alias
	}
	key := *channelID + ":" + commandName
	cachedCompiledTemplate, exists := template.templateMap[key]
	if exists {
		if cachedCompiledTemplate != nil {
			return cachedCompiledTemplate, true
		}
		return nil, false
	}
	result, error := GetChannelTemplate(channelID, &commandName)
	if error != nil {
		template.templateMap[key] = nil
		return nil, false
	}
	if result.AliasTo == result.CommandName {
		template.aliases[result.CommandName] = result.AliasTo
		key = result.AliasTo + ":" + commandName
	}
	log.Println(result.Template)
	if result.Template == "" {
		template.templateMap[key] = nil
		return nil, false
	}
	dbTemplate, templateError := mustache.ParseString(result.Template)
	if templateError != nil {
		return nil, false
	}
	template.templateMap[key] = dbTemplate
	return dbTemplate, true
}

//&chatMessage.User, &chatMessage.UserID, &chatMessage.ChannelID, &commandName, &commandName, &template
func (template templateContainer) UpdateTemplate(user *string, userID *string, channelID *string, commandName *string, templateBody *string) error {
	if *templateBody == "" {
		template.templateMap[*channelID+":"+*commandName] = nil

	} else {
		compiledTemplate, templateError := mustache.ParseString(*templateBody)
		if templateError != nil {
			//	log.Println(templateError)
			return templateError
		}
		template.templateMap[*channelID+":"+*commandName] = compiledTemplate

	}
	PutChannelTemplate(user, userID, channelID, commandName, commandName, templateBody)
	PushCommandsForChannel(channelID)
	return nil
}

func (template templateContainer) SetAliasto(user *string, userID *string, channelID *string, commandName *string, aliasTo *string) {
	template.aliases[*channelID+":"+*commandName] = *aliasTo
	aliasTemplate := ""

	result, error := GetChannelTemplate(channelID, aliasTo)
	if error == nil {
		aliasTemplate = result.Template
	}

	PutChannelTemplate(user, userID, channelID, commandName, aliasTo, &aliasTemplate)
	PushCommandsForChannel(channelID)

}

// TemplateCache is object, that updates and compiles templates which stored in memory, being backed up with database
var TemplateCache = templateContainer{make(map[string]*mustache.Template), make(map[string]string)}
