package commandHandlers

import (
	"log"

	"github.com/hoisie/mustache"
	"github.com/khades/servbot/repos"
)

type templateContainer struct {
	templateMap map[string]*mustache.Template
}

// Template is service DUUH
var Template = templateContainer{make(map[string]*mustache.Template)}

func (template templateContainer) get(channel string, commandName string) (bool, *mustache.Template) {
	key := channel + ":" + commandName
	cachedCompiledTemplate, exists := template.templateMap[key]
	if exists {
		return true, cachedCompiledTemplate
	}

	result, error := repos.GetChannelTemplate(channel, commandName)

	if error != nil {
		log.Print(error)
		template.templateMap[key] = nil
		return false, nil
	}

	if result.Template == "" {
		template.templateMap[key] = nil
		return false, nil
	}

	dbTemplate, templateError := mustache.ParseString(result.Template)

	if templateError != nil {
		log.Print(error)
		return false, nil

	}

	template.templateMap[key] = dbTemplate
	return true, dbTemplate

}

func (template templateContainer) update(channel string, commandName string, templateBody string) error {
	compiledTemplate, templateError := mustache.ParseString(templateBody)

	if templateError != nil {
		log.Println(templateError)
		return templateError
	}

	templates, error := repos.GetChannelAliasedTemplates(channel, commandName)

	if error != nil {
		log.Panic(error)
	}

	for _, item := range templates {
		key := channel + ":" + item.CommandName
		if templateBody == "" {
			template.templateMap[key] = nil
		}
		template.templateMap[key] = compiledTemplate
	}
	return nil
}
