package httpbackend

import (
	"log"
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func mod(next sessionAndChannelHandlerFunc) sessionAndChannelHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, session *models.HTTPSession, channelID *string, channelName *string) {
		channelInfo, error := repos.GetChannelInfo(channelID)
		log.Println(channelInfo)
		if error != nil {
			log.Println(error)
			writeJSONError(w, "That channel is not defined", http.StatusForbidden)
			return
		}
		if *channelID == session.UserID || channelInfo.GetIfUserIsMod(&session.UserID) == true {
			next(w, r, session, channelID, channelName)
		} else {
			writeJSONError(w, "You're not moderator", http.StatusForbidden)
			return
		}
	}
}

func withMod(next sessionAndChannelHandlerFunc) http.HandlerFunc {
	return withSessionAndChannel(mod(next))
}
