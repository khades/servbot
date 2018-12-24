package template

//TemplateAliasSelector is query selector for specific channel and command aliases
type TemplateAliasSelector struct {
	ChannelID string
	AliasTo   string
}

//TemplateSelector is query selector for specific channel and command without aliases
type TemplateSelector struct {
	ChannelID   string
	CommandName string
}
