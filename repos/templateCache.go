package repos

import (
	"log"

	"github.com/cbroglie/mustache"
)

type templateContainer struct {
	templateMap map[string]*mustache.Template
	aliases     map[string]string
}

func (template templateContainer) Get(channel *string, inputCommandName *string) (*mustache.Template, bool) {
	commandName := *inputCommandName
	alias, exists := template.aliases[*inputCommandName]
	if exists {
		commandName = alias
	}
	key := *channel + ":" + commandName
	cachedCompiledTemplate, exists := template.templateMap[key]
	if exists {
		if cachedCompiledTemplate != nil {
			return cachedCompiledTemplate, true
		}
		return nil, false
	}
	result, error := GetChannelTemplate(channel, &commandName)
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

func (template templateContainer) UpdateTemplate(channel *string, commandName *string, templateBody *string) error {
	if *templateBody == "" {
		template.templateMap[*channel+":"+*commandName] = nil
		return nil
	}
	compiledTemplate, templateError := mustache.ParseString(*templateBody)
	if templateError != nil {
		//	log.Println(templateError)
		return templateError
	}
	template.templateMap[*channel+":"+*commandName] = compiledTemplate
	return nil
}

func (template templateContainer) SetAliasto(channel *string, commandName *string, aliasTo *string) {
	template.aliases[*channel+":"+*commandName] = *aliasTo
}

// TemplateCache is object, that updates and compiles templates which stored in memory, being backed up with database
var TemplateCache = templateContainer{make(map[string]*mustache.Template), make(map[string]string)}
