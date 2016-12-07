package httpbackend

import (
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/khades/servbot/models"
)

func auth(next sessionHandlerFunc) sessionHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, s *models.HTTPSession) {
		log.Println(s)
		_, err := govalidator.ValidateStruct(s)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next(w, r, s)
	}
}

func withAuth(next sessionHandlerFunc) http.HandlerFunc {
	return withSession(auth(next))
}
