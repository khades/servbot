package httpSession

import (
	"github.com/khades/servbot/twitchAPIClient"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
)


type Service struct{
	collection      *mgo.Collection
	twitchAPIClient *twitchAPIClient.TwitchAPIClient
}

// Get returns information of user specified by his oauth key
func (service *Service) Get(oauthKey *string) (*HTTPSession, error) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "oauth",
		"action":  "Get"})
	result := httpSessionDBstruct{}
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"createdat": time.Now()}},
		ReturnNew: true,
	}

	_, error := service.collection.Find(bson.M{"key": *oauthKey}).Apply(change, &result)

	if error == nil {
		logger.Debug("User found in database")
		return &result.HTTPSession, nil
	}

	// TODO: getting user
	user, error := service.twitchAPIClient.GetUserByOauth(*oauthKey)
	if error == nil {
		logger.Debug("User not found on twitch")
		return  nil, error
	}

	result = httpSessionDBstruct{HTTPSession{
		Username: strings.ToLower(user.DisplayName),
		UserID:   user.ID,
		Key:      *oauthKey, 
		AvatarURL: user.ProfileImage},
		time.Now()}
	service.collection.Upsert(bson.M{"key": *oauthKey}, result)
	return &result.HTTPSession, nil
}
