package repos

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/khades/servbot/models"
)

type twitchUserRepsonse struct {
	Data []models.TwitchUserInfo `json:"data"`
}

func twitchHelix(method string, urlStr string, body io.Reader) (*http.Response, error) {
	var timeout = 5 * time.Second
	var client = http.Client{Timeout: timeout}
	req, error := http.NewRequest(method, "https://api.twitch.tv/helix/"+urlStr, body)
	req.Header.Add("Authorization", "Bearer "+Config.OauthKey)
	req.Header.Add("Client-ID", Config.ClientID)
	if error != nil {
		return nil, error
	}
	return client.Do(req)
}
func getUsersByParameterLimited(idSlice []string, idType string) ([]models.TwitchUserInfo, error) {
	var delimiter = "&" + idType + "="

	usersString := "users?" + idType + "=" + strings.Join(idSlice, delimiter)
	resp, error := twitchHelix("GET", usersString, nil)
	if error != nil {
		return nil, error
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	var twitchUserResponseStruct twitchUserRepsonse
	marshallError := json.NewDecoder(resp.Body).Decode(&twitchUserResponseStruct)
	if marshallError != nil {
		return nil, marshallError
	}
	return twitchUserResponseStruct.Data, nil
}
func getUsersByParameter(idList []string, idType string) ([]models.TwitchUserInfo, error) {

	var result []models.TwitchUserInfo
	sliceStart := 0
	for index := range idList {
		if index == len(idList)-1 || index-sliceStart > 48 {
			pageResult, error := getUsersByParameterLimited(idList[sliceStart:index+1], idType)

			if error != nil {
				return nil, error
			}

			result = append(result, pageResult...)
		}
	}
	return result, nil
}

func getTwitchUsersByDisplayName(displayNames []string) ([]models.TwitchUserInfo, error) {
	return getUsersByParameter(displayNames, "login")
}

func getTwitchUsersByID(userIDS []string) ([]models.TwitchUserInfo, error) {
	return getUsersByParameter(userIDS, "id")
}
