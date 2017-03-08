package httpbackend

import (
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type sessionHandlerFunc func(w http.ResponseWriter, r *http.Request, s *models.HTTPSession)

func session(next sessionHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := repos.GetSession(r)
		if err != nil {
			writeJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		val := session.Values["sessions"]
		var sessionObject = &models.HTTPSession{}
		if val != nil {
			var ok = false
			sessionObject, ok = val.(*models.HTTPSession)
			if ok == false {
				writeJSONError(w, "what", http.StatusInternalServerError)
				return
			}
		} else {
			if repos.Config.Debug == true {
				sessionObject = &models.HTTPSession{
					Username:  "khadesru",
					Key:       "123",
					UserID:    "40635840",
					AvatarURL: "https://static-cdn.jtvnw.net/jtv_user_pictures/khadesru-profile_image-d6f7260bc68376d7-300x300.jpeg"}
			}

		}
		next(w, r, sessionObject)
	}
}

func withSession(next sessionHandlerFunc) http.HandlerFunc {
	return corsEnabled(session(next))
}
