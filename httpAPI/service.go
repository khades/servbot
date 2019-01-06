package httpAPI

import (
	"github.com/khades/servbot/metrics"
	"net/http"
	"sync"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/gorilla/websocket"

	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/config"
	"goji.io"
	"goji.io/pat"

	"github.com/khades/servbot/httpSession"

	"github.com/sirupsen/logrus"
	"gopkg.in/asaskevich/govalidator.v4"
)

type Service struct {
	config             *config.Config
	httpSessionService *httpSession.Service
	channelInfoService *channelInfo.Service
	eventBus           EventBus.Bus
	metrics            *metrics.Service
	requestsCounter    map[string]requestCounterRecord
	mux                *goji.Mux
}

func (service *Service) NewMux() *goji.Mux {
	return service.mux
}

func (service *Service) Serve(wg *sync.WaitGroup) {
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		http.ListenAndServe("localhost:8000", service.mux)
		wg.Done()
	}(wg)
}

func (service *Service) auth(next SessionHandlerFunc) SessionHandlerFunc {

	return func(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession) {
		var authLogger = logrus.WithFields(logrus.Fields{
			"package": "httpAPI",
			"action":  "auth"})

		_, err := govalidator.ValidateStruct(s)
		// time.Sleep(2 * time.Second)
		if err != nil {
			WriteJSONError(w, "Not authorized", http.StatusUnauthorized)
			return
		}
		authLogger.Debugf("Incoming authorized request for userID %s", s.UserID)
		requestCounterForUser, found := service.requestsCounter[s.Key]
		if found == true && requestCounterForUser.Date.After(time.Now()) {
			service.requestsCounter[s.Key] = requestCounterRecord{Count: requestCounterForUser.Count + 1, Date: requestCounterForUser.Date}
		} else {
			service.requestsCounter[s.Key] = requestCounterRecord{Count: 1, Date: time.Now().Add(time.Minute)}
		}
		authLogger.Debugf("Current request count for user %s is %d", s.UserID, service.requestsCounter[s.Key].Count)
		service.metrics.LogHTTPApiUserRequest(s.Username)

		if service.requestsCounter[s.Key].Count > 200 {
			authLogger.Debugf("Rejecting user api request for userID %s ", s.UserID)

			WriteJSONError(w, "Too Many Requests", http.StatusTooManyRequests)

		} else {
			next(w, r, s)
		}

	}
}

func (service *Service) WithAuth(next SessionHandlerFunc) http.HandlerFunc {
	return service.WithSession(service.auth(next))
}

func (service *Service) session(next SessionHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if service.config.Debug == true {
			next(w, r, &httpSession.HTTPSession{
				AvatarURL: "https://static-cdn.jtvnw.net/jtv_user_pictures/d7bbfe70-9962-4e76-8a10-267a0cea9a64-profile_image-300x300.png",
				Username:  "khadesru",
				UserID:    "40635840",
				Key:       "fdfsfd"})
			return
		}
		cookie, err := r.Cookie("oauth")
		if err != nil || cookie.Value == "" {
			WriteJSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userInfo, userInfoError := service.httpSessionService.Get(&cookie.Value)
		if userInfoError != nil {
			WriteJSONError(w, userInfoError.Error(), http.StatusUnauthorized)
			return
		}

		next(w, r, userInfo)
	}
}

func (service *Service) WithSession(next SessionHandlerFunc) http.HandlerFunc {
	return service.corsEnabled(service.session(next))
}

func (service *Service) sessionAndChannel(next SessionAndChannelHandlerFunc) SessionHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession) {
		channelID := pat.Param(r, "channel")

		if channelID == "" {
			WriteJSONError(w, "channel variable is not defined", http.StatusUnprocessableEntity)
			return
		}

		channel, error := service.channelInfoService.Get(&channelID)

		if error != nil {
			WriteJSONError(w, error.Error(), http.StatusInternalServerError)
			return
		}
		service.metrics.LogHTTPApiChannelRequest(channel.Channel)

		next(w, r, s, channel)

	}
}

func (service *Service) WithSessionAndChannel(next SessionAndChannelHandlerFunc) http.HandlerFunc {
	return service.corsEnabled(service.session(service.auth(service.sessionAndChannel(next))))
}

func (service *Service) corsEnabled(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if service.config.Debug == true {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Methods", "POST, PUT, GET, OPTIONS")
			w.Header().Add("Access-Control-Allow-Headers", "X-Requested-With, Content-Type")
		}
		next(w, r)
	}
}

func owner(next SessionAndChannelHandlerFunc) SessionAndChannelHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, session *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
		if session.UserID == channelInfo.ChannelID {
			next(w, r, session, channelInfo)
		} else {
			WriteJSONError(w, "You're not owner", http.StatusForbidden)
			return
		}
	}
}

func mod(next SessionAndChannelHandlerFunc) SessionAndChannelHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, session *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {

		if channelInfo.GetIfUserIsMod(&session.UserID) == true {
			next(w, r, session, channelInfo)
		} else {
			WriteJSONError(w, "You're not moderator", http.StatusForbidden)
			return
		}
	}
}

func (service *Service) WithMod(next SessionAndChannelHandlerFunc) http.HandlerFunc {
	return service.WithSessionAndChannel(mod(next))
}

func (service *Service) WithOwner(next SessionAndChannelHandlerFunc) http.HandlerFunc {
	return service.WithSessionAndChannel(owner(next))
}

func (service *Service) CorsEnabled(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if service.config.Debug == true {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Methods", "POST, PUT, GET, OPTIONS")
			w.Header().Add("Access-Control-Allow-Headers", "X-Requested-With, Content-Type")
		}
		next(w, r)
	}
}

func (service *Service) Options(w http.ResponseWriter, r *http.Request) {
	service.CorsEnabled(Options)(w, r)
}
func (service *Service) WSEvent(w http.ResponseWriter, r *http.Request, messageID string) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "httpAPI",
		"action":  "WSEvent:" + messageID})
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true }}

	conn, err := upgrader.Upgrade(w, r, nil)
	pongWait := 40 * time.Second

	if err != nil {
		logger.Debug(err)
		return
	}

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	ping := func(value string) {
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			return
		}
	}
	service.eventBus.Subscribe("ping", ping)
	defer service.eventBus.Unsubscribe("ping", ping)

	write := func(value string) {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(value)); err != nil {
			return
		}
	}

	service.eventBus.Subscribe(messageID, write)
	defer service.eventBus.Unsubscribe(messageID, write)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			break
		}
	}
}
