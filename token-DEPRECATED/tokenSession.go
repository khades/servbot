package token_DEPRECATED

import (
	"net/http"

	"github.com/khades/servbot/repos"
)

type tokensessionFunc func(w http.ResponseWriter, r *http.Request, channelID *string)

func tokenSession(next tokensessionFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("oauth")

		if err != nil || cookie.Value == "" {
			writeJSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		token := r.URL.Query().Get("token")

		if token == "" {
			writeJSONError(w, "Token is not specified", http.StatusUnauthorized)
			return
		}

		channelID, err := repos.GetChannelToken(token)

		if err != nil {
			writeJSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		next(w, r, &channelID)
	}
}

func withTokenSession(next tokensessionFunc)   http.HandlerFunc {
	return corsEnabled(tokenSession(next))

}
