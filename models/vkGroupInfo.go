package models

type VkGroupInfo struct {
	GroupName       string
	NotifyOnChange  bool
	LastMessageID   int
	LastMessageURL  string
	LastMessageBody string
}
