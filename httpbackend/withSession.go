package httpbackend

import (
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type sessionHandlerFunc func(w http.ResponseWriter, r *http.Request, s *models.HTTPSession)

func session(next sessionHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if repos.Config.Debug == true {
			next(w, r, &models.HTTPSession{
				AvatarURL: "https://static-cdn.jtvnw.net/jtv_user_pictures/d7bbfe70-9962-4e76-8a10-267a0cea9a64-profile_image-300x300.png",
				Username:  "khadesru",
				UserID:    "40635840",
				Key:       "fdfsfd"})
			return
		}
		cookie, err := r.Cookie("oauth")
		if err != nil || cookie.Value == "" {
			writeJSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userInfo, userInfoError := repos.GetUserInfoByOauth(&cookie.Value)
		if userInfoError != nil {
			writeJSONError(w, userInfoError.Error(), http.StatusUnauthorized)
			return
		}

		next(w, r, userInfo)
	}
}

func withSession(next sessionHandlerFunc) http.HandlerFunc {
	return corsEnabled(session(next))
}
