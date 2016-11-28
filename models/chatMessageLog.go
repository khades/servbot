package models

type ChatMessageLog struct {
	User     string
	Channel  string
	Messages []MessageStruct
}
