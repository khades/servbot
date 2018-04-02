package repos

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/khades/servbot/models"
)

var httpsessionCollection = "httpSessions"

//db.log_events.createIndex( { "createdAt": 1 }, { expireAfterSeconds: 3600 } )
type httpSessionDBstruct struct {
	models.HTTPSession `bson:",inline"`
	CreatedAt          time.Time
}

// GetUserInfoByOauth returns information of user specified by his oauth key
func GetUserInfoByOauth(oauthKey *string) (*models.HTTPSession, error) {
	result := httpSessionDBstruct{}
	error := db.C(httpsessionCollection).Find(bson.M{"key": *oauthKey}).One(&result)
	if error != nil {

		return &result.HTTPSession, nil
	}

	nameResp, err := twitchHelixOauth("GET", "users", nil, *oauthKey)
	if err != nil {
		return nil, err
	}
	if nameResp != nil {
		defer nameResp.Body.Close()
	}
	if nameResp.StatusCode == 400 {
		return nil, errors.New("Twitch Error, Cant get username")

	}
	var usernameStruct = twitchUserRepsonse{}

	nameMarshallError := json.NewDecoder(nameResp.Body).Decode(usernameStruct)
	if nameMarshallError != nil {
		return nil, errors.New("Twitch Error, Cant marshall username: " + nameMarshallError.Error())
	}
	if len(usernameStruct.Data) == 0 {
		return nil, errors.New("No user found")
	}
	result = httpSessionDBstruct{models.HTTPSession{
		Username: strings.ToLower(usernameStruct.Data[0].DisplayName),
		UserID:   usernameStruct.Data[0].ID,
		Key:      *oauthKey, AvatarURL: usernameStruct.Data[0].ProfileImage},
		time.Now()}
	db.C(httpsessionCollection).Update(bson.M{"key": *oauthKey}, httpSessionDBstruct{})
	return &result.HTTPSession, nil
}
