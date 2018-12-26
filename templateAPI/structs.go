package templateAPI

import "github.com/khades/servbot/template"

type templatesResponse struct {
	Templates []template.TemplateInfo `json:"templates"`
	IsMod     bool                  `json:"isMod"`
}


type aliasToRequest struct {
	AliasTo string `json:"aliasTo"`
}