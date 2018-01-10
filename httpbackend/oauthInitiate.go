package httpbackend

import (
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func oauthInitiate(w http.ResponseWriter, r *http.Request) {
	cookie, cookieErr := r.Cookie("oauth")
	redirect := false
	if cookieErr != nil || cookie.Value == "" {
		redirect = true
	}

	userInfo, userInfoError := repos.GetUserInfoByOauth(&cookie.Value)
	if (userInfoError != nil) {
		redirect = true

	}

	_, err := govalidator.ValidateStruct(userInfo)
	if err != nil {
		redirect = true

	}
	if redirect == true {
		log.Println("redirecting to api")
		http.Redirect(w, r, "https://api.twitch.tv/kraken/oauth2/authorize"+
			"?response_type=code"+
			"&client_id="+repos.Config.ClientID+
			"&redirect_uri="+repos.Config.AppOauthURL+
			"&scope=user_subscriptions+user_read", http.StatusFound)

	} else {
		http.Redirect(w, r, repos.Config.AppURL+"/#/afterAuth", http.StatusFound)

	}


}
