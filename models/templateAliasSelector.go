package models

//TemplateAliasSelector is query selector for specific channel and command aliases
type TemplateAliasSelector struct {
	Channel string
	AliasTo string
}
