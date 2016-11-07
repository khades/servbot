package httpbackend

import (
	"fmt"
	"net/http"

	"github.com/khades/servbot/repos"
	"github.com/kidstuff/mongostore"

	goji "goji.io"

	"goji.io/pat"
)

var sessionStore = mongostore.NewMongoStore(repos.Db.C("sessions"), 3600, true, []byte("something-very-secret"))

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
func hello(w http.ResponseWriter, r *http.Request) {
	name := pat.Param(r, "name")
	fmt.Fprintf(w, "Hello, %s!", name)
}

// Start We are starting server here
func Start() {
	mux := goji.NewMux()
	mux.Handle(pat.New("/logs/*"), logs)
	mux.HandleFunc(pat.Get("/hello/:name"), hello)
	http.ListenAndServe("localhost:8000", mux)
}