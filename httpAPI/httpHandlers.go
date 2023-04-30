package httpAPI

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/httpSession"
	"github.com/sirupsen/logrus"
)

type tokenResponse struct {
	Token string `json:"access_token"`
}

func (service *Service) oauth(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "httpAPI",
		"action":  "oauth"})
	code := r.URL.Query().Get("code")
	if code == "" {
		WriteJSONError(w, "Incoming Twitch code is missing", http.StatusUnprocessableEntity)
		return
	}
	postValues := url.Values{
		"client_id":     {service.config.ClientID},
		"client_secret": {service.config.ClientSecret},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {service.config.AppOauthURL},
		"code":          {code}}

	// https://id.twitch.tv/oauth2/token
	resp, err := http.PostForm("https://id.twitch.tv/oauth2/token'", postValues)

	if err != nil {
		logger.Infof("Twitch Error: %s", err.Error())
		WriteJSONError(w, "Twitch Error, Cant get auth token, Connection problem", http.StatusUnprocessableEntity)
		return
	}

	if resp.StatusCode == 400 {
		body, err := ioutil.ReadAll(resp.Body)

		if err == nil {
			logger.Infof("Parsed body of first 400 error: %s", string(body))
		} else {
			logger.Infof("We didnt parsed body of first 400 error: %s", err.Error())
		}
		WriteJSONError(w, "Twitch Error, Cant get auth token, Got code 400", http.StatusUnprocessableEntity)
		return
	}
	bs := string(resp.Body)
	logger.Infof(bs)
	var tokenStruct = new(tokenResponse)

	marshallError := json.NewDecoder(resp.Body).Decode(tokenStruct)
	if marshallError != nil {
		logger.Infof("Decoding error: %s", marshallError.Error())
		WriteJSONError(w, "Twitch Error, Can't marshall oauth token", http.StatusUnprocessableEntity)
		return
	}
	expiration := time.Now().Add(3 * 24 * time.Hour)

	cookie := http.Cookie{Name: "oauth", Value: tokenStruct.Token, Expires: expiration}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, service.config.AppURL+"/#/afterAuth", http.StatusFound)
	defer resp.Body.Close()
}

func (service *Service) oauthInitiate(w http.ResponseWriter, r *http.Request) {
	isValidCookie := validateSession(r, service.httpSessionService)
	if isValidCookie == false {

		//"https://id.twitch.tv/oauth2/authorize"
		http.Redirect(w, r, "https://id.twitch.tv/oauth2/authorize"+
			"?response_type=code"+
			"&client_id="+service.config.ClientID+
			"&redirect_uri="+service.config.AppOauthURL+
			"&scope=user_subscriptions+user_read", http.StatusFound)

	} else {
		http.Redirect(w, r, service.config.AppURL+"/#/afterAuth", http.StatusFound)

	}

}
func (service *Service) userIndex(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession) {
	response := userIndexResponse{}
	if (s.UserID == service.config.BotUserID || s.UserID == "40635840") {
		channels, _ := service.channelInfoService.GetModChannelsForAdmin()
		response.ModChannels = channels
	}else {
		channels, _ := service.channelInfoService.GetModChannels(&s.UserID)
		response.ModChannels = channels
	}
	response.Username = s.Username
	response.AvatarURL = s.AvatarURL
	response.UserID = s.UserID
	json.NewEncoder(w).Encode(response)
}

func channelName(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	json.NewEncoder(w).Encode(channelNameStruct{channelInfo.Channel})
}

func channelInfoHandler(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	json.NewEncoder(w).Encode(&channelInfo)
}

func getTime(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(timeResponse{time.Now()})
}
