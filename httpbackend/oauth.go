package httpbackend

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type tokenResponse struct {
	Token string `json:"access_token"`
}
type nameResponse struct {
	Name string
}

func oauth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Incoming Twitch code is missing", http.StatusUnprocessableEntity)
		return
	}
	resp, err := http.PostForm("https://api.twitch.tv/kraken/oauth2/token",
		url.Values{
			"client_id":     {repos.Config.ClientID},
			"client_secret": {repos.Config.ClientSecret},
			"grant_type":    {"authorization_code"},
			"redirect_uri":  {repos.Config.AppOauthURL},
			"code":          {code}})
	//"state":         {}
	if err != nil {
		log.Println(err)
		http.Error(w, "Twitch Error", http.StatusUnprocessableEntity)
		return
	}
	var tokenStruct = new(tokenResponse)

	marshallError := json.NewDecoder(resp.Body).Decode(tokenStruct)
	if marshallError != nil {
		log.Println(marshallError)
		http.Error(w, "Twitch Error", http.StatusUnprocessableEntity)
		return
	}
	nameResp, err := http.Get("https://api.twitch.tv/kraken/user?client_id=" + repos.Config.ClientID + "&oauth_token=" + tokenStruct.Token)
	if err != nil {
		log.Println(err)
		http.Error(w, "Twitch Error", http.StatusUnprocessableEntity)
		return
	}

	var usernameStruct = new(nameResponse)
	nameMarshallError := json.NewDecoder(nameResp.Body).Decode(usernameStruct)
	if nameMarshallError != nil {
		log.Println(marshallError)
		http.Error(w, "Twitch Error", http.StatusUnprocessableEntity)
		return
	}
	log.Println("We got credentials of ", usernameStruct.Name)
	log.Println(nameResp.Body)

	session, err := repos.GetSession(r)
	session.Values["sessions"] = models.HTTPSession{Username: usernameStruct.Name, Key: tokenStruct.Token}
	session.Save(r, w)
	fmt.Fprintf(w, "Hello, %s!", usernameStruct.Name)

	defer resp.Body.Close()
}
