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

func twitchHelixOauth(method string, urlStr string, body io.Reader, key string) (*http.Response, error) {
	var timeout = 5 * time.Second
	var client = http.Client{Timeout: timeout}
	req, error := http.NewRequest(method, "https://api.twitch.tv/helix/"+urlStr, body)
	req.Header.Add("Authorization", "Bearer "+key)
	req.Header.Add("Client-ID", Config.ClientID)
	if error != nil {
		return nil, error
	}
	return client.Do(req)
}

func twitchHelix(method string, urlStr string, body io.Reader) (*http.Response, error) {
	return twitchHelixOauth(method, urlStr, body, Config.OauthKey)
}

func getUsersByParameterPaged(idSlice []string, idType string) ([]models.TwitchUserInfo, error) {
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

	streamsString := "streams?user_id=" + strings.Join(channels, ",user_id=")
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
