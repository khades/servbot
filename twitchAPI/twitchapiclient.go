package twitchAPI

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/khades/servbot/metrics"
	"io"
	"net/http"
	"net/http/httputil"

	"strings"
	"time"

	"github.com/khades/servbot/config"
	"github.com/sirupsen/logrus"
)

type Client struct {
	config *config.Config
	metrics *metrics.Service
}

func Init(config *config.Config, metrics *metrics.Service) *Client {
	return &Client{
		config: config,
		metrics: metrics,
	}
}

func (tApi *Client) twitchHelixPost(urlStr string, body io.Reader) (*http.Response, error) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "twitchAPI",
		"action":  "twitchHelixPost"})
	var timeout = 5 * time.Second
	var client = http.Client{Timeout: timeout}
	logger.Debugf("Url: %s", urlStr)
	req, error := http.NewRequest("POST", "https://api.twitch.tv/helix/"+urlStr, body)
	req.Header.Add("Authorization", "Bearer "+strings.Replace(tApi.config.APIKey, "oauth:", "", 1))
	req.Header.Add("Client-ID", tApi.config.ClientID)
	req.Header.Add("Content-Type", "application/json")
	dump, dumpErr := httputil.DumpRequestOut(req, true)
	if dumpErr == nil {
		logger.Debugf("Request is %q", dump)
	}
	if error != nil {
		return nil, error
	}
	tApi.metrics.LogTwitchRequest()
	return client.Do(req)
}

func (tApi *Client) twitchHelixOauth(method string, urlStr string, body io.Reader, key string) (*http.Response, error) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "twitchAPI",
		"action":  "twitchHelixOauth"})
	var timeout = 5 * time.Second
	var client = http.Client{Timeout: timeout}
	logger.Debugf("Url: %s", urlStr)
	req, error := http.NewRequest(method, "https://api.twitch.tv/helix/"+urlStr, body)
	req.Header.Add("Authorization", "Bearer "+strings.Replace(key, "oauth:", "", 1))
	req.Header.Add("Client-ID", tApi.config.ClientID)

	if error != nil {
		return nil, error
	}

	tApi.metrics.LogTwitchRequest()
	res, err := client.Do(req)

	if err != nil {
		logger.Infof("Error: %s", err.Error())
		return nil, err
	}

	dump, dumpErr := httputil.DumpResponse(res, true)
	if dumpErr == nil {
		logger.Debugf("Repsonse is %q", dump)
	}

	return res, err
}

func (tApi *Client) twitchHelix(method string, urlStr string, body io.Reader) (*http.Response, error) {
	return tApi.twitchHelixOauth(method, urlStr, body, tApi.config.APIKey)
}

// func (tApi *Client) getFollowers(channelID *string, noCursor bool) (*twitchFollowerResponse, error) {
// 	logger := logrus.WithFields(logrus.Fields{
// 		"package": "repos",
// 		"feature": "twitchapi",
// 		"action":  "getFollowers"})
// 	url := "users/follows?to_id=" + *channelID
// 	if noCursor == true {
// 		url = url + "&first=1"

// 	}
// 	logger.Debugf("Request url is %s", url)

// 	resp, error := tApi.twitchHelix("GET", url, nil)
// 	if error != nil {
// 		return nil, error
// 	}
// 	if resp != nil {
// 		defer resp.Body.Close()
// 	}
// 	// dump, err := httputil.DumpResponse(resp, true)
// 	// if err != nil {

// 	// 	log.Fatal(err)
// 	// }
// 	// logger.Debugf("Repsonse is %q", dump)

// 	var twitchResponseStruct twitchFollowerResponse
// 	marshallError := json.NewDecoder(resp.Body).Decode(&twitchResponseStruct)
// 	if marshallError != nil {
// 		logger.Debugf("Marshalling error: %s", marshallError.Error())
// 		return nil, marshallError
// 	}
// 	return &twitchResponseStruct, nil
// }

func (tApi *Client) SubscribeToChannelFollowerWebhook(channelID string, secret string) bool {
	logger := logrus.WithFields(logrus.Fields{
		"package": "twitchAPI",
		"action":  "SubscribeToChannelFollowerWebhook"})
	form := hub{
		Mode:         "subscribe",
		Topic:        "https://api.twitch.tv/helix/users/follows?to_id=" + channelID,
		Callback:     "https://servbot.khades.org/api/webhook/follows",
		LeaseSeconds: "864000",
		Secret:       secret}

	body, _ := json.Marshal(form)

	tApi.metrics.LogTwitchSpecificRequest("SubscribeToChannelFollowerWebhook")
	resp, _ := tApi.twitchHelixPost("webhooks/hub", bytes.NewReader(body))
	if resp != nil {
		defer resp.Body.Close()
	}
	dump, err := httputil.DumpResponse(resp, true)
	if err == nil {
		logger.Debugf("Repsonse is %q", dump)
	} else {
		return false
	}
	logger.Debugf("Status is %d", resp.StatusCode)
	if resp.StatusCode == http.StatusAccepted {
		return true
	}
	return false
}

func (tApi *Client) GetUserFollowDate(channelID *string, userID *string) (bool, time.Time) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "twitchAPI",
		"action":  "GetUserFollowDate"})
	url := "users/follows?from_id=" + *userID + "&to_id=" + *channelID

	logger.Debugf("Request url is %s", url)

	tApi.metrics.LogTwitchSpecificRequest("GetUserFollowDate")
	resp, error := tApi.twitchHelix("GET", url, nil)
	if error != nil {
		return false, time.Now()
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	var twitchResponseStruct twitchFollowerResponse
	marshallError := json.NewDecoder(resp.Body).Decode(&twitchResponseStruct)
	if marshallError != nil {
		logger.Debugf("Marshalling error: %s", marshallError.Error())
		return false, time.Now()
	}
	logger.Debugf("Response %+v", twitchResponseStruct)

	logger.Debugf("Checking if user %s follows channel %s %s: %s", *userID, *channelID, len(twitchResponseStruct.Followers), twitchResponseStruct.Total == 0)
	if len(twitchResponseStruct.Followers) == 0 {
		return false, time.Now()
	}
	return true, twitchResponseStruct.Followers[0].Date
}

func (tApi *Client) getUsersByParameterPaged(idSlice []string, idType string) ([]TwitchUserInfo, error) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "twitchAPI",
		"action":  "getUsersByParameterPaged"})
	var delimiter = "&" + idType + "="

	usersString := "users?" + idType + "=" + strings.Join(idSlice, delimiter)
	logger.Debugf("Url string to get users from twitch: %s", usersString)

	tApi.metrics.LogTwitchSpecificRequest("getUsersByParameterPaged")
	resp, error := tApi.twitchHelix("GET", usersString, nil)
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

func (tApi *Client) getUsersByParameter(idList []string, idType string) ([]TwitchUserInfo, error) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "twitchAPI",
		"action":  "getUsersByParameter"})
	logger.Debugf("Fetching users: %s", strings.Join(idList, ", "))
	var result []TwitchUserInfo
	sliceStart := 0
	for index := range idList {

		if index == len(idList)-1 || index-sliceStart > 48 {
			pageResult, error := tApi.getUsersByParameterPaged(idList[sliceStart:index+1], idType)

			if error != nil {
				return nil, error
			}

			result = append(result, pageResult...)
		}
	}
	return result, nil
}

func (tApi *Client) GetTwitchUsersByDisplayName(displayNames []string) ([]TwitchUserInfo, error) {
	return tApi.getUsersByParameter(displayNames, "login")
}

func (tApi *Client) GetTwitchUsersByID(userIDS []string) ([]TwitchUserInfo, error) {
	return tApi.getUsersByParameter(userIDS, "id")
}

func (tApi *Client) getStreamStatusesPaged(channels []string) ([]TwitchStreamStatus, error) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "twitchAPI",
		"action":  "getStreamStatusesPaged"})
	var twitchResult twitchStreamResponse

	streamsString := "streams?user_id=" + strings.Join(channels, "&user_id=")
	logger.Debugf("Url: %s", streamsString)

	tApi.metrics.LogTwitchSpecificRequest("getStreamStatusesPaged")
	resp, error := tApi.twitchHelix("GET", streamsString, nil)
	if error != nil {
		logger.Infof("Error: %s", error.Error())
		return nil, error
	}
	dumpedResponse, _ := httputil.DumpResponse(resp, true)
	logger.Debugf("Response: %s", string(dumpedResponse[:]))
	if resp != nil {
		defer resp.Body.Close()
	}
	marshallError := json.NewDecoder(resp.Body).Decode(&twitchResult)
	if marshallError != nil {
		return nil, marshallError
	}
	return twitchResult.Data, nil
}

// getStreamStatuses returns stream information for all active channels for that chatbot instance
func (tApi *Client) GetStreamStatuses() ([]TwitchStreamStatus, error) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "twitchAPI",
		"action":  "getStreamStatuses"})
	var twitchResult []TwitchStreamStatus
	sliceStart := 0

	for index := range tApi.config.ChannelIDs {
		if index == len(tApi.config.ChannelIDs)-1 || index-sliceStart > 48 {
			pageResult, error := tApi.getStreamStatusesPaged(tApi.config.ChannelIDs[sliceStart : index+1])
			if error != nil {
				logger.Infof("Api error: %+v", error)
				return nil, error
			}
			twitchResult = append(twitchResult, pageResult...)
		}
	}
	return twitchResult, nil
}

func (tApi *Client) getGamesByIDPaged(gamesIDS []string) ([]twitchGameStruct, error) {
	var twitchResult twitchGameStructResponse

	gamesString := "games?id=" + strings.Join(gamesIDS, "&id=")

	tApi.metrics.LogTwitchSpecificRequest("getGamesByIDPaged")
	resp, error := tApi.twitchHelix("GET", gamesString, nil)
	if error != nil {
		return nil, error
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	marshallError := json.NewDecoder(resp.Body).Decode(&twitchResult)
	if marshallError != nil {
		return nil, marshallError
	}
	return twitchResult.Data, nil
}

func (tApi *Client) GetGamesByID(gameIDS []string) ([]twitchGameStruct, error) {
	var twitchResult []twitchGameStruct
	sliceStart := 0
	for index := range gameIDS {
		if index == len(gameIDS)-1 || index-sliceStart > 48 {
			pageResult, error := tApi.getGamesByIDPaged(gameIDS[sliceStart : index+1])

			if error != nil {
				return nil, error
			}

			twitchResult = append(twitchResult, pageResult...)
		}
	}
	return twitchResult, nil
}

func (tApi *Client) GetUserByOauth(key string) (*TwitchUserInfo, error) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "twitchAPI",
		"action":  "GetUserByOauth"})

	tApi.metrics.LogTwitchSpecificRequest("GetUserByOauth")

	nameResp, err := tApi.twitchHelixOauth("GET", "users", nil, key)
	if err != nil {
		logger.Debug("Twitch helix error : %s", err.Error())

		return nil, err
	}

	if nameResp != nil {
		defer nameResp.Body.Close()
	}
	if nameResp.StatusCode == 400 {
		logger.Info("Twitch Error, Cant get username")

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
	return &usernameStruct.Data[0], nil
}

func (tApi *Client) GetAPIKey() (*ApiKeyResponse, error) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "twitchAPI",
		"action":  "GetAPIKey"})
	var timeout = 5 * time.Second
	var client = http.Client{Timeout: timeout}
	url := "https://id.twitch.tv/oauth2/token?client_id=" + tApi.config.ClientID + "&client_secret=" + tApi.config.ClientSecret + "&grant_type=client_credentials"
	logger.Debugf("Url: %s", url)
	req, reqError := http.NewRequest("POST", url, nil)
	if reqError != nil {
		logger.Debug("Twitch Error, General error: " + reqError.Error())

		return nil, reqError
	}
	tApi.metrics.LogTwitchRequest()
	tApi.metrics.LogTwitchSpecificRequest("GetAPIKey")

	resp, error := client.Do(req)
	if error == nil {
		defer resp.Body.Close()
	} else {
		logger.Debug("Twitch Error, General error: " + error.Error())
		return nil, error
	}
	var apiKey = twitchAPIKeyResponse{}

	marshallError := json.NewDecoder(resp.Body).Decode(&apiKey)
	if marshallError != nil {
		logger.Debug("Twitch Error, Cant marshall username: " + marshallError.Error())

		return nil, errors.New("Twitch Error, Cant marshall response")
	}
	expirationDate := time.Now().Add(time.Duration(apiKey.ExpiresIn) * time.Second)
	logger.Debug("Result: %+v", apiKey)
	return &ApiKeyResponse{
		AccessToken:  apiKey.AccessToken,
		ExpiresIn:    expirationDate,
	}, nil
}

// func (tApi *Client) RefreshAPIKey() {
// 	apiUrl := "https://id.twitch.tv/oauth2/token"
// 	data := url.Values{}
// 	data.Set("grant_type", "refresh_token")
// 	data.Set("refresh_token", tApi.config.APIKeyRefreshToken)
// 	data.Set("client_id", tApi.config.ClientID)
// 	data.Set("client_secret", tApi.config.ClientSecret)

// 	client := &http.Client{}
// 	r, _ := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode())) // URL-encoded payload
// 	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

// 	resp, _ := client.Do(r)
// }
