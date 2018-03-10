package models

//ChatMessageLog struct describes specific user logs on specific channel
type ChatMessageLog struct {
	User           string          `json:"user"`
	Channel        string          `json:"channel"`
	KnownNicknames []string        `json:"knownNicknames"`
	UserID         string          `json:"userID"`
	ChannelID      string          `json:"channelID"`
	Messages       []MessageStruct `json:"messages"`
	Bans           []BanInfo       `json:"bans"`
}
