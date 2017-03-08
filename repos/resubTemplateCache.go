package repos

import (
	"errors"

	"github.com/hoisie/mustache"
	"github.com/khades/servbot/models"
)

type resubTemplateContainer map[string]*mustache.Template

// ResubTemplateCache is needed to cache mustache templates for resub message
var ResubTemplateCache resubTemplateContainer = make(map[string]*mustache.Template)

func (container resubTemplateContainer) put(channelID *string, template *mustache.Template) {
	container[*channelID] = template
}

func (container resubTemplateContainer) Get(subAlert *models.SubAlert) (*mustache.Template, error) {
	if subAlert.ResubMessage == "" {
		return nil, errors.New("empty string")
	}
	template, error := mustache.ParseString(subAlert.ResubMessage)
	if error != nil {
		return nil, error
	}
	container[subAlert.ChannelID] = template
	return template, nil
}

func (container resubTemplateContainer) Drop(channelID *string) {
	container[*channelID] = nil
}
