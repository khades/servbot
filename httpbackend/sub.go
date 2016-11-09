package httpbackend

import (
	"net/http"

	"goji.io/pat"

	"github.com/khades/servbot/models"
)

func sub(next sessionHandler) sessionHandler {
	return func(w http.ResponseWriter, r *http.Request, session *models.HTTPSession) {
		channel := pat.Param(r, "channel")
		if channel == "" {
			http.Error(w, "channel is not defined", http.StatusUnprocessableEntity)
			return
		}

		next(w, r, session)
	}
}
