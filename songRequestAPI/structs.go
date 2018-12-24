package songRequestAPI

import "github.com/khades/servbot/songRequest"

type songRequestResponse struct {
	*songRequest.ChannelSongRequest
	IsMod   bool `json:"isMod"`
	IsOwner bool `json:"isOwner"`
	// Token string `json:"token"`
}
