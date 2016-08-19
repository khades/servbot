package models

// TemplateInfo describes info about chat template WITHOUT history
type TemplateInfo struct {
	CommandName string
	*ChannelSelector
	AliasTo  string
	Template string
}
