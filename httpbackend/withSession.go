package httpbackend

import (
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type sessionHandler func(w http.ResponseWriter, r *http.Request, s *models.HTTPSession)

func withSession(next sessionHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := repos.GetSession(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		val := session.Values["sessions"]
		var sessionObject = &models.HTTPSession{}
		sessionObject, ok := val.(*models.HTTPSession)
		if ok == false {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
			// Handle the case that it's not an expected type
		}

		next(w, r, sessionObject)
	})
}
