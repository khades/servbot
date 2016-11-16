package httpbackend

import (
	"net/http"

	"goji.io/pat"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func mod(next sessionHandlerFunc) sessionHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, session *models.HTTPSession) {
		channel := pat.Param(r, "channel")
		if channel == "" {
			http.Error(w, "channel is not defined", http.StatusUnprocessableEntity)
			return
		}
		channelInfo, error := repos.GetChannelInfo(&channel)
		if error != nil {
			http.Error(w, error.Error(), http.StatusInternalServerError)
			return
		}
		isMod := false
		for _, value := range channelInfo.Mods {
			if value == session.Username {
				isMod = true
				break
			}
		}
		if isMod == true {
			next(w, r, session)
		} else {
			http.Error(w, "You're not moderator", http.StatusForbidden)
			return
		}
	}
}
