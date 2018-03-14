package httpbackend

import (
	"net/http"

	"goji.io/pat"

	"github.com/khades/servbot/httpclient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func sub(next sessionHandlerFunc) sessionHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, session *models.HTTPSession) {
		channel := pat.Param(r, "channel")
		if channel == "" {
			writeJSONError(w, "channel is not defined", http.StatusUnprocessableEntity)
			return
		}
		isSub, found := repos.GetIfSubToChannel(&session.Username, &channel)
		if found == false {
			url := "https://api.twitch.tv/kraken/users/" + session.Username + "/subscriptions/" + channel + "?oauth_token=" + session.Key
			resp, respError := httpclient.Get(url)
			if respError == nil && (resp.StatusCode == 200 || resp.StatusCode == 204) {
				isSub = true
			}
			if resp != nil {
				defer resp.Body.Close()
			}
			repos.SetIfSubToChannel(&session.Username, &channel, &isSub)
		}
		if isSub == true {
			next(w, r, session)
		} else {
			writeJSONError(w, "You're not sub", http.StatusForbidden)
		}
	}
}

func withSub(next sessionHandlerFunc) http.HandlerFunc {
	return withAuth(sub(next))
}
