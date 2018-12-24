package channelLogs

import (
	"time"

	"github.com/khades/servbot/chatMessage"
)

//ChatMessageLog struct describes specific user logs on specific channel
type ChatMessageLog struct {
	User           string                      `json:"user"`
	Channel        string                      `json:"channel"`
	KnownNicknames []string                    `json:"knownNicknames"`
	UserID         string                      `json:"userID"`
	ChannelID      string                      `json:"channelID"`
	Messages       []chatMessage.MessageStruct `json:"messages"`
	Bans           []BanInfo                   `json:"bans"`
}

// BanInfo struct describes info about person being banned
type BanInfo struct {
	User        string    `json:"user"`
	Duration    int       `json:"duration"`
	Reason      string    `json:"reason"`
	BanIssuer   string    `json:"banIssuer"`
	BanIssuerID string    `json:"banIssuerID"`
	Type        string    `json:"type"`
	Date        time.Time `json:"date"`
}
