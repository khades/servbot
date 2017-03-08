package httpbackend

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/khades/servbot/httpclient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
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
		log.Println(err)
		writeJSONError(w, "Twitch Error, Cant get auth token, Connection problem", http.StatusUnprocessableEntity)
		return
	}

	if resp.StatusCode == 400 {
		body, err := ioutil.ReadAll(resp.Body)

		if err == nil {

			log.Println(string(body))
		} else {
			log.Println("We didnt parsed body of first 400 error")

			log.Println(err)
		}
		writeJSONError(w, "Twitch Error, Cant get auth token, Got code 400", http.StatusUnprocessableEntity)
		return
	}
	var tokenStruct = new(tokenResponse)

	marshallError := json.NewDecoder(resp.Body).Decode(tokenStruct)
	if marshallError != nil {
		log.Println(marshallError)
		writeJSONError(w, "Twitch Error, Can't marshall oauth token", http.StatusUnprocessableEntity)
		return
	}
	url := "https://api.twitch.tv/kraken/users/" + repos.Config.ClientID
	nameResp, err := httpclient.TwitchV5(repos.Config.ClientID, "GET", url, nil)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "Twitch Error, Cant get username", http.StatusUnprocessableEntity)
		return
	}

	if nameResp.StatusCode == 400 {
		body, err := ioutil.ReadAll(resp.Body)

		if err == nil {
			log.Println(string(body))
		} else {
			log.Println(err)
			log.Println("We didnt parsed body of username request")
		}
		writeJSONError(w, "Twitch Error, Cant get username", http.StatusUnprocessableEntity)
		return
	}
	var usernameStruct = new(nameResponse)

	nameMarshallError := json.NewDecoder(nameResp.Body).Decode(usernameStruct)
	if nameMarshallError != nil {
		log.Println(marshallError)
		writeJSONError(w, "Twitch Error, Cant marshall username", http.StatusUnprocessableEntity)
		return
	}
	session, err := repos.GetSession(r)
	session.Options.Path = "/"
	sessionObject := models.HTTPSession{Username: usernameStruct.Name, UserID: usernameStruct.ID, Key: tokenStruct.Token, AvatarURL: usernameStruct.Logo}
	session.Values["sessions"] = sessionObject
	log.Println(sessionObject)
	session.Save(r, w)
	http.Redirect(w, r, repos.Config.AppURL+"/#/afterAuth", http.StatusFound)
	defer resp.Body.Close()
}
