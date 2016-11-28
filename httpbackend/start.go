package httpbackend

import (
	"net/http"

	"github.com/khades/servbot/repos"
	"github.com/kidstuff/mongostore"

	goji "goji.io"

	"goji.io/pat"
)

var sessionStore = mongostore.NewMongoStore(repos.Db.C("sessions"), 3600, true, []byte("something-very-secret"))

// Start We are starting server here
func Start() {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/api/channel/:channel/logs"), corsEnabled(withSession(logsUsers)))
	mux.HandleFunc(pat.Get("/api/channel/:channel/logs/:user"), corsEnabled(withSession(logs)))
	mux.HandleFunc(pat.Get("/api/user"), corsEnabled(user))
	mux.HandleFunc(pat.Get("/api/channel/:channel/templates"), corsEnabled(withSession(templates)))
	mux.HandleFunc(pat.Get("/api/channel/:channel/templates/:template"), corsEnabled(withSession(template)))

	// mux.HandleFunc(pat.Get("/isMod/:channel"), withSession(mod(func(w http.ResponseWriter, r *http.Request, session *models.HTTPSession) {
	// 	fmt.Fprintf(w, "Hello, %s, you're moderator of that channel!", session.Username)
	// })))
	// mux.HandleFunc(pat.Get("/isSub/:channel"), withSession(sub(func(w http.ResponseWriter, r *http.Request, session *models.HTTPSession) {
	// 	fmt.Fprintf(w, "Hello, %s, you're moderator of that channel!", session.Username)
	// })))
	mux.HandleFunc(pat.Get("/oauth"), oauth)
	mux.HandleFunc(pat.Get("/oauth/initiate"), withSession(oauthInitiate))
	http.ListenAndServe("localhost:8000", mux)
}
