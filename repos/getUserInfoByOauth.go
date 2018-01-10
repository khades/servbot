package repos

import (
	"errors"
	"encoding/json"

	"github.com/khades/servbot/httpclient"
	"github.com/khades/servbot/models"
)

type requestError struct {
	Error   string
	Status  string
	Message string
}
type tokenResponse struct {
	Token string `json:"access_token"`
}
type nameResponse struct {
	Name string `json:"name"`
	ID   string `json:"_id"`
	Logo string `json:"logo"`
}

func GetUserInfoByOauth(oauthKey *string) (*models.HTTPSession, error) {
	cacheKey := "OauthKey:"+*oauthKey
	binaryData, found := cacheObject.Get(cacheKey)
	if found {
		result :=binaryData.(models.HTTPSession)
		return  &result, nil
	}
	url := "https://api.twitch.tv/kraken/user?oauth_token=" + *oauthKey
	nameResp, err := httpclient.TwitchV5(Config.ClientID, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	if nameResp.StatusCode == 400 {
		return nil, errors.New("Twitch Error, Cant get username")

	}
	var usernameStruct = new(nameResponse)

	nameMarshallError := json.NewDecoder(nameResp.Body).Decode(usernameStruct)
	if nameMarshallError != nil {
		return nil, errors.New("Twitch Error, Cant marshall username: "+nameMarshallError.Error())
	}

    result :=  models.HTTPSession{Username: usernameStruct.Name, UserID: usernameStruct.ID, Key: *oauthKey, AvatarURL: usernameStruct.Logo}
	return &result, nil
}
