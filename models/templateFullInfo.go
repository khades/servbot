package models

// TemplateFullInfo describes info about template with history
type TemplateFullInfo struct {
	*TemplateInfo
	History []TemplateHistoryItem
}
