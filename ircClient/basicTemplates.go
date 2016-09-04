package ircClient

import (
	"log"

	"github.com/hoisie/mustache"
)

type basicTemplates struct {
	PublicTemplate            *mustache.Template
	PublicNonTargetedTemplate *mustache.Template
	PrivateTemplate           *mustache.Template
}

func generateTemplates() basicTemplates {
	publicTemplate, publicTemplateError := mustache.ParseString("PRIVMSG #{{ Channel }} :@{{ User }} {{Body}}")
	if publicTemplateError != nil {
		log.Panicln(publicTemplateError)
	}
	publicNonTargetedTemplate, publicNonTargetedTemplateError := mustache.ParseString("PRIVMSG #{{ Channel }} :{{Body}}")
	if publicNonTargetedTemplateError != nil {
		log.Panicln(publicNonTargetedTemplateError)
	}
	privateTemplate, privateTemplateError := mustache.ParseString("PRIVMSG #jtv /w {{ User }} Channel #{{ Channel }} as:{{ Body }}")
	if privateTemplateError != nil {
		log.Panicln(privateTemplateError)
	}
	return basicTemplates{publicTemplate, publicNonTargetedTemplate, privateTemplate}
}

var basicTemplatesInstance = generateTemplates()
