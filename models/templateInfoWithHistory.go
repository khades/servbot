package models

type TemplateInfoWithHistory struct {
	TemplateInfo `bson:",inline"`
	History      []TemplateHistory `json:"history"`
}
