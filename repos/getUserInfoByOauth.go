package repos

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/khades/servbot/models"
	"github.com/sirupsen/logrus"
)

var httpsessionCollection = "httpSessions"

//db.log_events.createIndex( { "createdAt": 1 }, { expireAfterSeconds: 3600 } )
type httpSessionDBstruct struct {
	models.HTTPSession `bson:",inline"`
	CreatedAt          time.Time
}

// GetUserInfoByOauth returns information of user specified by his oauth key
func GetUserInfoByOauth(oauthKey *string) (*models.HTTPSession, error) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "oauth",
		"action":  "GetUserInfoByOauth"})
	result := httpSessionDBstruct{}
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"createdat": time.Now()}},
		ReturnNew: true,
	}

	_, error := db.C(httpsessionCollection).Find(bson.M{"key": *oauthKey}).Apply(change, &result)

	if error == nil {
		logger.Debug("User found in database")
		return &result.HTTPSession, nil
	}

	nameResp, err := twitchHelixOauth("GET", "users", nil, *oauthKey)
	if err != nil {
		logger.Debug("Twitch helix error : %s", err.Error())

		return nil, err
	}

	if nameResp != nil {
		defer nameResp.Body.Close()
	}
	if nameResp.StatusCode == 400 {
		logger.Debug("Twitch Error, Cant get username")

		return nil, errors.New("Twitch Error, Cant get username")

	}
	var usernameStruct = twitchUserRepsonse{}

	nameMarshallError := json.NewDecoder(nameResp.Body).Decode(&usernameStruct)
	if nameMarshallError != nil {
		logger.Debug("Twitch Error, Cant marshall username: " + nameMarshallError.Error())

		return nil, errors.New("Twitch Error, Cant marshall username: " + nameMarshallError.Error())
	}
	if len(usernameStruct.Data) == 0 {
		logger.Debug("No user found")
		return nil, errors.New("No user found")
	}
	result = httpSessionDBstruct{models.HTTPSession{
		Username: strings.ToLower(usernameStruct.Data[0].DisplayName),
		UserID:   usernameStruct.Data[0].ID,
		Key:      *oauthKey, 
		AvatarURL: usernameStruct.Data[0].ProfileImage},
		time.Now()}
	db.C(httpsessionCollection).Upsert(bson.M{"key": *oauthKey}, result)
	return &result.HTTPSession, nil
}
