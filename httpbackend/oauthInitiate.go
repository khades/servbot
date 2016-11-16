package httpbackend

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func oauthInitiate(w http.ResponseWriter, r *http.Request, s *models.HTTPSession) {
	_, err := govalidator.ValidateStruct(s)
	if err != nil {
		http.Redirect(w, r, "https://api.twitch.tv/kraken/oauth2/authorize"+
			"?response_type=code"+
			"&client_id="+repos.Config.ClientID+
			"&redirect_uri="+repos.Config.AppOauthURL+
			"&scope=user_subscriptions+user_read", http.StatusFound)
		return
	}
	http.Redirect(w, r, repos.Config.AppOauthURL+"/#/afterAuth", http.StatusFound)
	return

}