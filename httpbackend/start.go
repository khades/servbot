package httpbackend

import (
	"net/http"

	goji "goji.io"

	"goji.io/pat"
)

// Start We are starting server here
func Start() {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/api/channel/:channel"), withSessionAndChannel(channel))
	mux.HandleFunc(pat.Get("/api/channel/:channel/channelname"), withSessionAndChannel(channelName))
	mux.HandleFunc(pat.Get("/api/channel/:channel/logs"), withMod(logsUsers))
	mux.HandleFunc(pat.Get("/api/channel/:channel/logs/search/:search"), withMod(logsUsersSearch))

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
	// mux.HandleFunc(pat.Get("/api/channel/:channel/bits"), withMod(bits))
	// mux.HandleFunc(pat.Get("/api/channel/:channel/bits/search/:search"), withMod(bitsSearch))
	// mux.HandleFunc(pat.Get("/api/channel/:channel/bits/:userID"), withMod(userbits))

	mux.HandleFunc(pat.Get("/api/channel/:channel/templates"), withSessionAndChannel(templates))

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

	mux.HandleFunc(pat.Get("/api/channel/:channel/automessages/removeinactive"), withMod(autoMessageRemoveInactive))
	mux.HandleFunc(pat.Options("/api/channel/:channel/automessages/removeinactive"), corsEnabled(options))

	mux.HandleFunc(pat.Get("/api/channel/:channel/automessages/:messageID"), withMod(autoMessageGet))
	mux.HandleFunc(pat.Options("/api/channel/:channel/automessages/:id"), corsEnabled(options))

	mux.HandleFunc(pat.Post("/api/channel/:channel/automessages"), withMod(autoMessageCreate))
	mux.HandleFunc(pat.Post("/api/channel/:channel/automessages/:id"), withMod(autoMessageUpdate))

	mux.HandleFunc(pat.Get("/oauth"), oauth)
	mux.HandleFunc(pat.Get("/oauth/initiate"), oauthInitiate)

	mux.HandleFunc(pat.Get("/api/time"), corsEnabled(getTime))
	mux.HandleFunc(pat.Get("/api/timeticker"), corsEnabled(timeTicker))

	mux.HandleFunc(pat.Get("/api/channel/:channel/subtrain"), withSessionAndChannel(subtrain))
	mux.HandleFunc(pat.Options("/api/channel/:channel/subtrain"), corsEnabled(options))
	mux.HandleFunc(pat.Post("/api/channel/:channel/subtrain"), withMod(putSubtrain))
	mux.HandleFunc(pat.Get("/api/channel/:channel/subtrain/events"), withMod(subtrainEvents))

	mux.HandleFunc(pat.Get("/api/channel/:channel/bans"), withMod(channelBans))
	mux.HandleFunc(pat.Options("/api/channel/:channel/bans"), corsEnabled(options))

	mux.HandleFunc(pat.Get("/api/channel/:channel/subdays"), withMod(subdayList))
	mux.HandleFunc(pat.Options("/api/channel/:channel/subdays"), corsEnabled(options))

	mux.HandleFunc(pat.Post("/api/channel/:channel/subdays/new"), withMod(subdayCreate))
	mux.HandleFunc(pat.Options("/api/channel/:channel/subdays/new"), corsEnabled(options))

	mux.HandleFunc(pat.Get("/api/channel/:channel/subdays/:subdayID"), withSessionAndChannel(subdayByID))
	mux.HandleFunc(pat.Options("/api/channel/:channel/subdays/:subdayID"), corsEnabled(options))

	mux.HandleFunc(pat.Get("/api/channel/:channel/subdays/:subdayID/close"), withMod(subdayClose))
	mux.HandleFunc(pat.Options("/api/channel/:channel/subdays/:subdayID/close"), corsEnabled(options))

	mux.HandleFunc(pat.Get("/api/channel/:channel/subdays/:subdayID/randomize"), withMod(subdayRandomize))
	mux.HandleFunc(pat.Options("/api/channel/:channel/subdays/:subdayID/randomize"), corsEnabled(options))

	mux.HandleFunc(pat.Get("/api/channel/:channel/subdays/:subdayID/pullwinner/:user"), withMod(subdayPullWinner))
	mux.HandleFunc(pat.Options("/api/channel/:channel/subdays/:subdayID/pullwinner/:user"), corsEnabled(options))

	mux.HandleFunc(pat.Get("/api/channel/:channel/songrequests"), withSessionAndChannel(songrequests))

	mux.HandleFunc(pat.Options("/api/channel/:channel/songrequests"), corsEnabled(options))

	mux.HandleFunc(pat.Get("/api/channel/:channel/songrequests/skip/:videoID"), withSessionAndChannel(songrequestsSkip))
	mux.HandleFunc(pat.Get("/api/channel/:channel/songrequests/bubbleup/:videoID"), withSessionAndChannel(songrequestsBubbleUp))
	mux.HandleFunc(pat.Get("/api/channel/:channel/songrequests/bubbleuptosecond/:videoID"), withSessionAndChannel(songrequestsBubbleUpToSecond))

	mux.HandleFunc(pat.Get("/api/channel/:channel/songrequests/events"), withSessionAndChannel(songrequestsEvents))
	mux.HandleFunc(pat.Post("/api/channel/:channel/songrequests/settings"), withMod(songrequestsPushSettings))
	mux.HandleFunc(pat.Options("/api/channel/:channel/songrequests/settings"), corsEnabled(options))

	mux.HandleFunc(pat.Get("/api/channel/:channel/songrequests/setvolume/:volume"), withMod(songrequestSetVolume))
	mux.HandleFunc(pat.Options("/api/channel/:channel/songrequests/setvolume/:volume"), corsEnabled(options))

	mux.HandleFunc(pat.Get("/api/channel/:channel/songrequests/:videoID/settag/:tag"), withMod(songrequestSetTag))
	mux.HandleFunc(pat.Options("/api/channel/:channel/songrequests/:videoID/settag/:tag"), corsEnabled(options))

	mux.HandleFunc(pat.Get("/api/channel/:channel/songrequests/:videoID/unban"), withMod(songrequestsUnban))
	mux.HandleFunc(pat.Options("/api/channel/:channel/songrequests/:videoID/unban"), corsEnabled(options))

	mux.HandleFunc(pat.Get("/api/channel/:channel/songrequests/library/get"), withMod(songrequestGetLibrary))
	mux.HandleFunc(pat.Get("/api/channel/:channel/songrequests/bannedtracks"), withMod(songrequestGetBannedTracks))

	mux.HandleFunc(pat.Get("/api/webhook/streams"), corsEnabled(webhookVerify))
	mux.HandleFunc(pat.Post("/api/webhook/streams"), corsEnabled(webhookStream))
	mux.HandleFunc(pat.Options("/api/webhook/streams"), corsEnabled(options))

	mux.HandleFunc(pat.Get("/api/webhook/follows"), corsEnabled(webhookVerify))
	mux.HandleFunc(pat.Post("/api/webhook/follows"), corsEnabled(webhookFollows))
	mux.HandleFunc(pat.Options("/api/webhook/follows"), corsEnabled(options))

	http.ListenAndServe("localhost:8000", mux)
}
