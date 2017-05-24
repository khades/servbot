package models

type VkGroupInfo struct {
	GroupName       string `json:"groupName"`
	NotifyOnChange  bool   `json:"notifyOnChange"`
	LastMessageID   int    `json:"lastMessageID"`
	LastMessageURL  string `json:"lastMessageURL"`
	LastMessageBody string `json:"lastMessageBody"`
}
