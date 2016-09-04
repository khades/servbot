package models

import "strings"

// ChannelInfo describes data of
type ChannelInfo struct {
	Channel      string
	StreamStatus StreamStatus
	Mods         []string
	Commands     []string
}

// Helper Command for mustashe
func (channelInfo ChannelInfo) GetCommands() string {
	return strings.Join(channelInfo.Commands, ", ")
}
