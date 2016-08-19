package models

//TemplateAliasSelector defines format of query from templates collection
type TemplateAliasSelector struct {
	*ChannelSelector
	AliasTo string
}
