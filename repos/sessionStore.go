package repos

import (
	"net/http"

	"github.com/gorilla/sessions"
)

//mongostore "gopkg.in/go-playground/mongostore.v4"

//var sessionStore = mongostore.New(dbSession, "sessions", &sessions.Options{MaxAge: 3600, Secure: true}, true,
//[]byte("secret-key"))
var sessionStore = sessions.NewCookieStore([]byte("ASFasfafasfasfSfAS"))

// GetSession is configured object for reading http sessions
func GetSession(r *http.Request) (*sessions.Session, error) {
	return sessionStore.Get(r, "sessionkey")
}
