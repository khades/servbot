package templateAPI

import "github.com/khades/servbot/template"

type templatesResponse struct {
	Templates []template.TemplateInfo `json:"templates"`
	IsMod     bool                  `json:"isMod"`
}

type templatePushRequest struct {
	Template string `json:"get"`
}

type aliasToRequest struct {
	AliasTo string `json:"aliasTo"`
}