package httpbackend

import (
	"net/http"

	"goji.io/pat"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type sessionAndChannelHandlerFunc func(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channel *string)

func sessionAndChannel(next sessionAndChannelHandlerFunc) sessionHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, s *models.HTTPSession) {
		channelID := pat.Param(r, "channel")

		if channelID == "" {
			writeJSONError(w, "channel variable is not defined", http.StatusUnprocessableEntity)
			return
		}
		if channelID == s.UserID {
			next(w, r, s, &channelID, &s.Username)
			return

		}
		channel, error := repos.GetChannelInfo(&channelID)

		if error != nil {
			writeJSONError(w, error.Error(), http.StatusInternalServerError)
			return
		}
		next(w, r, s, &channelID, &channel.Channel)

	}
}

func withSessionAndChannel(next sessionAndChannelHandlerFunc) http.HandlerFunc {
	return corsEnabled(session(auth(sessionAndChannel(next))))
}
