package httpbackend

import (
	"log"
	"net/http"

	"goji.io/pat"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func mod(next sessionHandlerFunc) sessionHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, session *models.HTTPSession) {
		log.Println(session)
		channel := pat.Param(r, "channel")
		if channel == "" {
			writeJSONError(w, "Сhannel variable is not defined", http.StatusUnprocessableEntity)
			return
		}
		channelInfo, error := repos.GetChannelInfo(&channel)
		log.Println(channelInfo)
		if error != nil {
			log.Println(error)
			writeJSONError(w, "That channel is not defined", http.StatusForbidden)
			return
		}
		if channelInfo.GetIfUserIsMod(&session.Username) == true {
			next(w, r, session)
		} else {
			writeJSONError(w, "You're not moderator", http.StatusForbidden)
			return
		}
	}
}

func withMod(next sessionHandlerFunc) http.HandlerFunc {
	return withAuth(mod(next))
}
