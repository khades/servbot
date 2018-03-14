package repos

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/khades/servbot/models"
)

// GetUserInfoByOauth returns information of user specified by his oauth key
func GetUserInfoByOauth(oauthKey *string) (*models.HTTPSession, error) {
	cacheKey := "OauthKey:" + *oauthKey
	binaryData, found := cacheObject.Get(cacheKey)
	if found {
		result := binaryData.(models.HTTPSession)
		return &result, nil
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
	result := models.HTTPSession{Username: strings.ToLower(usernameStruct.Data[0].DisplayName), UserID: usernameStruct.Data[0].ID, Key: *oauthKey, AvatarURL: usernameStruct.Data[0].ProfileImage}
	return &result, nil
}
