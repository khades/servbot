package models

// TemplateInfo describes info about chat template WITHOUT history
type TemplateInfo struct {
	CommandName string
	Channel     string
	AliasTo     string
	Template    string
}
