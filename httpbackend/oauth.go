package httpbackend

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/khades/servbot/repos"
	"github.com/sirupsen/logrus"
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

func oauth(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "httpbackend",
		"feature": "oauth",
		"action":  "oauth"})
	code := r.URL.Query().Get("code")
	if code == "" {
		writeJSONError(w, "Incoming Twitch code is missing", http.StatusUnprocessableEntity)
		return
	}
	postValues := url.Values{
		"client_id":     {repos.Config.ClientID},
		"client_secret": {repos.Config.ClientSecret},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {repos.Config.AppOauthURL},
		"code":          {code}}
	resp, err := http.PostForm("https://api.twitch.tv/kraken/oauth2/token", postValues)

	if err != nil {
		logger.Infof("Twitch Error: %s",err.Error())
		writeJSONError(w, "Twitch Error, Cant get auth token, Connection problem", http.StatusUnprocessableEntity)
		return
	}

	if resp.StatusCode == 400 {
		body, err := ioutil.ReadAll(resp.Body)

		if err == nil {

			logger.Infof("Parsed body of first 400 error: %s", string(body))
		} else {
			logger.Infof("We didnt parsed body of first 400 error: %s", err.Error())
		}
		writeJSONError(w, "Twitch Error, Cant get auth token, Got code 400", http.StatusUnprocessableEntity)
		return
	}
	var tokenStruct = new(tokenResponse)

	marshallError := json.NewDecoder(resp.Body).Decode(tokenStruct)
	if marshallError != nil {
		logger.Infof("Decoding error: %s", marshallError.Error())
		writeJSONError(w, "Twitch Error, Can't marshall oauth token", http.StatusUnprocessableEntity)
		return
	}
	expiration := time.Now().Add(3 * 24 * time.Hour)

	cookie := http.Cookie{Name: "oauth", Value: tokenStruct.Token, Expires: expiration}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, repos.Config.AppURL+"/#/afterAuth", http.StatusFound)
	defer resp.Body.Close()
}
