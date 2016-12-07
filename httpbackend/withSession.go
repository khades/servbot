package httpbackend

import (
	"log"
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

		if val == nil {
			log.Println("session is nil")
			session.Values["sessions"] = models.HTTPSession{}
		} else {
			var ok = false
			log.Println(val)
			sessionObject, ok = val.(*models.HTTPSession)
			if ok == false {
				http.Error(w, "what", http.StatusInternalServerError)
				return
			}
		}
		log.Println("Returning session")
		log.Println(sessionObject)
		session.Save(r, w)
		next(w, r, sessionObject)
	}
}
