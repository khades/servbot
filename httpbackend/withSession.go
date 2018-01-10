package httpbackend

import (
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type sessionHandlerFunc func(w http.ResponseWriter, r *http.Request, s *models.HTTPSession)

func session(next sessionHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("oauth")
		if err != nil || cookie.Value == "" {
			writeJSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userInfo, userInfoError := repos.GetUserInfoByOauth(&cookie.Value)
		if (userInfoError != nil) {
			writeJSONError(w, userInfoError.Error(), http.StatusUnauthorized)
			return
		}
	
		next(w, r, userInfo)
	}
}

func withSession(next sessionHandlerFunc) http.HandlerFunc {
	return corsEnabled(session(next))
}
