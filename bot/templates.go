package bot

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
	var compiledTemplate *mustache.Template
	var found = false
	key := channel + ":" + commandName
	cachedCompiledTemplate, exists := template.templateMap[key]
	if exists {
		found = true
		compiledTemplate = cachedCompiledTemplate
	} else {
		result, error := repos.GetChannelTemplate(channel, commandName)
		if error != nil {
			log.Print(error)
		} else {
			dbTemplate, templateError := mustache.ParseString(result.Template)
			if templateError == nil {
				template.templateMap[key] = dbTemplate
				compiledTemplate = dbTemplate
			} else {
				log.Print(error)
			}
		}
	}
	return found, compiledTemplate
}

func (template templateContainer) update(channel string, commandName string, templateBody string, user string) {
	dbTemplate, templateError := mustache.ParseString(templateBody)
	if templateError == nil {
		key := channel + ":" + commandName
		template.templateMap[key] = dbTemplate
		repos.PutChannelTemplate(user, channel, commandName, templateBody)
	}
}
