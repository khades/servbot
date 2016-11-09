package httpbackend

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/khades/servbot/models"
)

func auth(next sessionHandler) sessionHandler {
	return func(w http.ResponseWriter, r *http.Request, s *models.HTTPSession) {
		_, err := govalidator.ValidateStruct(s)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
		}

		next(w, r, s)
	}
}
