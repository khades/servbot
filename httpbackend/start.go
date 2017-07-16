package httpbackend

import (
	"net/http"

	"github.com/khades/servbot/repos"
	"github.com/kidstuff/mongostore"

	goji "goji.io"

	"goji.io/pat"
)

var sessionStore = mongostore.NewMongoStore(repos.Db.C("sessions"), 3600*24, true, []byte("something-very-secret"))

// Start We are starting server here
func Start() {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/api/channel/:channel"), withSessionAndChannel(channel))
	mux.HandleFunc(pat.Get("/api/channel/:channel/channelname"), withSessionAndChannel(channelName))
	mux.HandleFunc(pat.Get("/api/channel/:channel/logs"), withMod(logsUsers))
	mux.HandleFunc(pat.Get("/api/channel/:channel/logs/username/:user"), withMod(logsByUsername))
	mux.HandleFunc(pat.Get("/api/channel/:channel/logs/userid/:userID"), withMod(logsByUserID))

	mux.HandleFunc(pat.Get("/api/channel/:channel/subs"), withMod(subscriptions))
	mux.HandleFunc(pat.Get("/api/channel/:channel/subs/limit/:limit"), withMod(subscriptionsWithLimit))
	mux.HandleFunc(pat.Get("/api/channel/:channel/info"), withMod(channelInfo))
	mux.HandleFunc(pat.Post("/api/channel/:channel/externalservices/vk"), withMod(putVK))
	mux.HandleFunc(pat.Options("/api/channel/:channel/externalservices/vk"), corsEnabled(options))
	mux.HandleFunc(pat.Post("/api/channel/:channel/externalservices/twitchdj"), withMod(putTwitchDJ))
	mux.HandleFunc(pat.Options("/api/channel/:channel/externalservices/twitchdj"), corsEnabled(options))
	mux.HandleFunc(pat.Get("/api/channel/:channel/subs/events"), withMod(subscriptionEvents))
	mux.HandleFunc(pat.Get("/api/user"), withAuth(user))
	mux.HandleFunc(pat.Get("/api/user/index"), withAuth(userIndex))
	mux.HandleFunc(pat.Get("/api/channel/:channel/bits"), withMod(bits))
	mux.HandleFunc(pat.Get("/api/channel/:channel/bits/:userID"), withMod(userbits))

	mux.HandleFunc(pat.Get("/api/channel/:channel/templates"), withMod(templates))

	mux.HandleFunc(pat.Get("/api/channel/:channel/templates/:commandName"), withMod(template))
	mux.HandleFunc(pat.Post("/api/channel/:channel/templates/:commandName"), withMod(putTemplate))
	mux.HandleFunc(pat.Options("/api/channel/:channel/templates/:commandName"), corsEnabled(options))

	mux.HandleFunc(pat.Post("/api/channel/:channel/templates/:commandName/setAliasTo"), withMod(aliasTemplate))
	mux.HandleFunc(pat.Options("/api/channel/:channel/templates/:commandName/setAliasTo"), corsEnabled(options))

	mux.HandleFunc(pat.Get("/api/channel/:channel/subalert"), withMod(subAlert))
	mux.HandleFunc(pat.Post("/api/channel/:channel/subalert"), withMod(setSubAlert))
	mux.HandleFunc(pat.Options("/api/channel/:channel/subalert"), corsEnabled(options))

	mux.HandleFunc(pat.Get("/api/channel/:channel/automessages"), withMod(autoMessageList))
	mux.HandleFunc(pat.Options("/api/channel/:channel/automessages"), corsEnabled(options))

	mux.HandleFunc(pat.Get("/api/channel/:channel/automessages/:messageID"), withMod(autoMessageGet))
	mux.HandleFunc(pat.Options("/api/channel/:channel/automessages/:id"), corsEnabled(options))

	mux.HandleFunc(pat.Post("/api/channel/:channel/automessages"), withMod(autoMessageCreate))
	mux.HandleFunc(pat.Post("/api/channel/:channel/automessages/:id"), withMod(autoMessageUpdate))

	mux.HandleFunc(pat.Get("/oauth"), oauth)
	mux.HandleFunc(pat.Get("/oauth/initiate"), withSession(oauthInitiate))

	mux.HandleFunc(pat.Get("/api/time"), corsEnabled(getTime))
	mux.HandleFunc(pat.Get("/api/timeticker"), corsEnabled(timeTicker))
	http.ListenAndServe("localhost:8000", mux)
}
