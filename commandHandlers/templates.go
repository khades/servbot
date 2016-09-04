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

func (template templateContainer) get(channel string, commandName string) *mustache.Template {
	key := channel + ":" + commandName
	cachedCompiledTemplate, exists := template.templateMap[key]
	if exists {
		log.Print("We found template")
		if cachedCompiledTemplate != nil {
			return cachedCompiledTemplate
		}
		return nil
	}

	result, error := repos.GetChannelTemplate(channel, commandName)

	if error != nil {

		log.Print(error)
		log.Print("Nothing's found")

		template.templateMap[key] = nil
		return nil
	}
	log.Println(result.Template)
	if result.Template == "" {
		log.Print("We found template, but it is empty")

		template.templateMap[key] = nil
		return nil
	}

	dbTemplate, templateError := mustache.ParseString(result.Template)

	if templateError != nil {
		log.Print(error)
		return nil

	}

	template.templateMap[key] = dbTemplate
	return dbTemplate

}
func (template templateContainer) updateTemplate(channel string, commandName string, templateBody string) error {
	if templateBody == "" {
		template.templateMap[channel+":"+commandName] = nil
		return nil
	}

	compiledTemplate, templateError := mustache.ParseString(templateBody)

	if templateError != nil {
		log.Println(templateError)
		return templateError
	}

	template.templateMap[channel+":"+commandName] = compiledTemplate

	return nil
}
func (template templateContainer) updateAliases(channel string, commandName string, templateBody string) error {
	compiledTemplate, templateError := mustache.ParseString(templateBody)

	if templateError != nil {
		log.Println(templateError)
		return templateError
	}
	templates, error := repos.GetChannelAliasedTemplates(channel, commandName)
	log.Println(templates)
	if error != nil {
		return error
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
