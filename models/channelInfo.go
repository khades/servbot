package models

// ChannelInfo describes data of
type ChannelInfo struct {
	Channel      string
	StreamStatus *StreamStatus
	Mods         []string
}
