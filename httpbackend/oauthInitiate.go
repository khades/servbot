package httpbackend

import (
	"net/http"

	"gopkg.in/asaskevich/govalidator.v4"
	"github.com/khades/servbot/repos"
)
func validateSession(r *http.Request) bool {
	cookie, cookieErr := r.Cookie("oauth")

	if cookieErr != nil || cookie == nil || cookie.Value == "" {
		return false
	} 

	userInfo, userInfoError := repos.GetUserInfoByOauth(&cookie.Value)
	if (userInfoError != nil) {
		return false

	}

	_, err := govalidator.ValidateStruct(userInfo)
	if err != nil {
		return false

	}
	return true
}

func oauthInitiate(w http.ResponseWriter, r *http.Request) {
	isValidCookie := validateSession(r)
	if isValidCookie == false {

		http.Redirect(w, r, "https://id.twitch.tv/oauth2/authorize"+
			"?response_type=code"+
			"&client_id="+repos.Config.ClientID+
			"&redirect_uri="+repos.Config.AppOauthURL+
			"&scope=user_subscriptions+user_read", http.StatusFound)

	} else {
		http.Redirect(w, r, repos.Config.AppURL+"/#/afterAuth", http.StatusFound)

	}


}
