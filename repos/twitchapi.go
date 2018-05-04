package repos

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

	"strings"
	"time"

	"github.com/khades/servbot/models"
	"github.com/sirupsen/logrus"
)

type twitchUserInfo struct {
	ID          string `json:"_id"`
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
}

type usersLoginStruct struct {
	Users []twitchUserInfo `json:"users"`
}

type twitchUserRepsonse struct {
	Data []models.TwitchUserInfo `json:"data"`
}
type twitchFollower struct {
	UserID    string    `json:"from_id"`
	ChannelID string    `json:"to_id"`
	Date      time.Time `json:"followed_at"`
}
type twitchFollowerResponsePagination struct {
	Cursor string `json:"cursor"`
}
type twitchFollowers []twitchFollower

func (followers twitchFollowers) Len() int {
	return len(followers)
}

func (followers twitchFollowers) Less(i, j int) bool {
	return followers[i].Date.Before(followers[j].Date)
}

func (followers twitchFollowers) Swap(i, j int) {
	followers[i], followers[j] = followers[j], followers[i]
}

type twitchFollowerResponse struct {
	Total      int64                            `json:"total"`
	Pagination twitchFollowerResponsePagination `json:"pagination"`
	Followers  twitchFollowers                  `json:"data"`
}

func twitchHelixPost(urlStr string, body io.Reader) (*http.Response, error) {
	var timeout = 5 * time.Second
	var client = http.Client{Timeout: timeout}
	req, error := http.NewRequest("POST", "https://api.twitch.tv/helix/"+urlStr, body)
	req.Header.Add("Authorization", "Bearer "+strings.Replace(Config.OauthKey, "oauth:", "", 1))
	req.Header.Add("Client-ID", Config.ClientID)
	req.Header.Add("Content-Type", "application/json")
	if error != nil {
		return nil, error
	}
	return client.Do(req)
}
func twitchHelixOauth(method string, urlStr string, body io.Reader, key string) (*http.Response, error) {
	var timeout = 5 * time.Second
	var client = http.Client{Timeout: timeout}
	req, error := http.NewRequest(method, "https://api.twitch.tv/helix/"+urlStr, body)
	req.Header.Add("Authorization", "Bearer "+strings.Replace(key, "oauth:", "", 1))
	req.Header.Add("Client-ID", Config.ClientID)

	if error != nil {
		return nil, error
	}
	return client.Do(req)
}

func getFollowers(channelID *string, noCursor bool) (*twitchFollowerResponse, error) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "followers",
		"action":  "getFollowers"})
	url := "users/follows?to_id=" + *channelID
	if noCursor == true {
		url = url + "&first=1"

	}
	logger.Debugf("Request url is %s", url)

	resp, error := twitchHelix("GET", url, nil)
	if error != nil {
		return nil, error
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {

		log.Fatal(err)
	}
	logger.Debugf("Repsonse is %q", dump)

	var twitchResponseStruct twitchFollowerResponse
	marshallError := json.NewDecoder(resp.Body).Decode(&twitchResponseStruct)
	if marshallError != nil {
		logger.Debugf("Marshalling error: %s", marshallError.Error())
		return nil, marshallError
	}
	return &twitchResponseStruct, nil
}

func twitchHelix(method string, urlStr string, body io.Reader) (*http.Response, error) {
	return twitchHelixOauth(method, urlStr, body, Config.OauthKey)
}
func getUserFollowDate(channelID *string, userID *string) (bool, time.Time) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "followers",
		"action":  "getUserFollowDate"})
	url := "users/follows?from_id=" + *userID + "&to_id=" + *channelID

	logger.Debugf("Request url is %s", url)

	resp, error := twitchHelix("GET", url, nil)
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
func getUsersByParameterPaged(idSlice []string, idType string) ([]models.TwitchUserInfo, error) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "followers",
		"action":  "getUsersByParameterPaged"})
	var delimiter = "&" + idType + "="

	usersString := "users?" + idType + "=" + strings.Join(idSlice, delimiter)
	logger.Debugf("Url string to get users from twitch: %s", usersString)
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
	logger := logrus.WithFields(logrus.Fields{
		"package": "repos",
		"feature": "followers",
		"action":  "getUsersByParameter"})
	logger.Debugf("Fetching users: %s", strings.Join(idList, ", "))
	var result []models.TwitchUserInfo
	sliceStart := 0
	for index := range idList {

		if index == len(idList)-1 || index-sliceStart > 48 {
			pageResult, error := getUsersByParameterPaged(idList[sliceStart:index+1], idType)

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

type twitchStreamResponse struct {
	Data []models.TwitchStreamStatus `json:"data"`
}

func getStreamStatusesPaged(channels []string) ([]models.TwitchStreamStatus, error) {
	var twitchResult twitchStreamResponse

	streamsString := "streams?user_id=" + strings.Join(channels, "&user_id=")
	resp, error := twitchHelix("GET", streamsString, nil)
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

// getStreamStatuses returns stream information for all active channels for that chatbot instance
func getStreamStatuses() ([]models.TwitchStreamStatus, error) {
	var twitchResult []models.TwitchStreamStatus
	sliceStart := 0

	for index := range Config.ChannelIDs {
		if index == len(Config.ChannelIDs)-1 || index-sliceStart > 48 {
			pageResult, error := getStreamStatusesPaged(Config.ChannelIDs[sliceStart : index+1])
			if error != nil {
				return nil, error
			}
			twitchResult = append(twitchResult, pageResult...)
		}
	}
	return twitchResult, nil
}

type twitchGameStruct struct {
	GameID string `json:"id"`
	Game   string `json:"name"`
}

type twitchGameStructResponse struct {
	Data []twitchGameStruct `json:"data"`
}

func getGamesByIDPaged(gamesIDS []string) ([]twitchGameStruct, error) {
	var twitchResult twitchGameStructResponse

	gamesString := "games?id=" + strings.Join(gamesIDS, ",id=")
	resp, error := twitchHelix("GET", gamesString, nil)
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

func getGamesByID(gameIDS []string) ([]twitchGameStruct, error) {
	var twitchResult []twitchGameStruct
	sliceStart := 0
	for index := range gameIDS {
		if index == len(gameIDS)-1 || index-sliceStart > 48 {
			pageResult, error := getGamesByIDPaged(gameIDS[sliceStart : index+1])

			if error != nil {
				return nil, error
			}

			twitchResult = append(twitchResult, pageResult...)
		}
	}
	return twitchResult, nil
}
