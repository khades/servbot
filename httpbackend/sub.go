package httpbackend

import (
	"log"
	"net/http"

	"goji.io/pat"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func sub(next sessionHandlerFunc) sessionHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, session *models.HTTPSession) {
		channel := pat.Param(r, "channel")
		if channel == "" {
			http.Error(w, "channel is not defined", http.StatusUnprocessableEntity)
			return
		}
		isSub, found := repos.GetIfSubToChannel(&session.Username, &channel)
		if found == false {
			log.Println("Sub cache is not set for user " + channel + " for channel " + session.Username)
			url := "https://api.twitch.tv/kraken/users/" + session.Username + "/subscriptions/" + channel + "?oauth_token=" + session.Key
			resp, respError := http.Get(url)
			if respError == nil && (resp.StatusCode == 200 || resp.StatusCode == 204) {
				isSub = true
			}
			repos.SetIfSubToChannel(&session.Username, &channel, &isSub)
			defer resp.Body.Close()
		}
		if isSub == true {
			next(w, r, session)
		} else {
			http.Error(w, "You're not sub", http.StatusForbidden)
		}
	}
}

func withSub(next sessionHandlerFunc) http.HandlerFunc {
	return withAuth(sub(next))
}
