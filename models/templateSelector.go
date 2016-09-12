package models

//TemplateSelector is query selector for specific channel and command without aliases
type TemplateSelector struct {
	Channel     string
	CommandName string
}
