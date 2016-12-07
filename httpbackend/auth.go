package httpbackend

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/khades/servbot/models"
)

func auth(next sessionHandlerFunc) sessionHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, s *models.HTTPSession) {
		_, err := govalidator.ValidateStruct(s)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"result":"","error":%q}`, err)
			json.NewEncoder(w).Encode(&models.HttpError{Code: 401, Message: "not authorized"})
			return
		}
		next(w, r, s)
	}
}

func withAuth(next sessionHandlerFunc) http.HandlerFunc {
	return withSession(auth(next))
}
