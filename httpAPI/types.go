package httpAPI

import (
	"net/http"
	"time"

	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/httpSession"
)

type SessionAndChannelHandlerFunc func(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelID *channelInfo.ChannelInfo)

type SessionHandlerFunc func(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession)

type requestCounterRecord struct {
	Count int
	Date  time.Time
}

type userIndexResponse struct {
	ModChannels []channelInfo.ChannelWithID `json:"modChannels"`
	Username    string                 `json:"username"`
	AvatarURL   string                 `json:"avatarUrl"`
	UserID      string                 `json:"userID"`
}

type channelNameStruct struct {
	Channel string `json:"channel"`
}

type OptionResponse struct {
	Status string
}

type timeResponse struct {
	Time time.Time `json:"time"`
}
