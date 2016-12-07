package repos

import (
	"net/http"

	mongostore "gopkg.in/go-playground/mongostore.v4"

	"github.com/gorilla/sessions"
)

var sessionStore = mongostore.New(dbSession, "sessions", &sessions.Options{MaxAge: 3600, Secure: true}, true,
	[]byte("secret-key"))

// GetSession is configured object for reading http sessions
func GetSession(r *http.Request) (*sessions.Session, error) {
	return sessionStore.Get(r, "sessionkey")
}
