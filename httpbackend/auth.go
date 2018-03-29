package httpbackend

import (
	"net/http"
	"time"

	"github.com/khades/servbot/models"
	"github.com/sirupsen/logrus"
	"gopkg.in/asaskevich/govalidator.v4"
)

type requestCounterRecord struct {
	Count int
	Date  time.Time
}

var authLogger = logrus.WithFields(logrus.Fields{
	"package": "httpbackend",
	"feature": "auth",
	"action":  "auth"})

var requestsCounter = make(map[string]requestCounterRecord)

func auth(next sessionHandlerFunc) sessionHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, s *models.HTTPSession) {
		_, err := govalidator.ValidateStruct(s)
		// time.Sleep(2 * time.Second)
		if err != nil {
			writeJSONError(w, "Not authorized", http.StatusUnauthorized)
			return
		}
		authLogger.Debugf("Incoming authorized request for userID %s", s.UserID)
		requestCounterForUser, found := requestsCounter[s.Key]
		if found == true && requestCounterForUser.Date.After(time.Now()) {
			requestsCounter[s.Key] = requestCounterRecord{Count: requestCounterForUser.Count + 1, Date: requestCounterForUser.Date}
		} else {
			requestsCounter[s.Key] = requestCounterRecord{Count: 1, Date: time.Now().Add(time.Minute)}
		}
		authLogger.Debugf("Current request count for user %s is %d", s.UserID, requestsCounter[s.Key].Count)

		if requestsCounter[s.Key].Count > 50 {
			authLogger.Debugf("Rejecting user api request for userID %s ", s.UserID)

			writeJSONError(w, "Too Many Requests", http.StatusTooManyRequests)

		} else {
			next(w, r, s)
		}

	}
}

func withAuth(next sessionHandlerFunc) http.HandlerFunc {
	return withSession(auth(next))
}
