package httpbackend

import (
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type sessionHandlerFunc func(w http.ResponseWriter, r *http.Request, s *models.HTTPSession)

func withSession(next sessionHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := repos.GetSession(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		val := session.Values["sessions"]
		var sessionObject = &models.HTTPSession{}

		if val != nil {
			var ok = false
			sessionObject, ok = val.(*models.HTTPSession)
			if ok == false {
				http.Error(w, "what", http.StatusInternalServerError)
				return
			}
		}
		next(w, r, sessionObject)
	}
}
