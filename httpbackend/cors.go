package httpbackend

import (
	"net/http"

	"github.com/khades/servbot/repos"
)

func corsEnabled(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if repos.Config.Debug == true {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Methods", "POST, PUT, GET, OPTIONS")
			w.Header().Add("Access-Control-Allow-Headers", "X-Requested-With, Content-Type")
		}
		next(w, r)
	}
}
