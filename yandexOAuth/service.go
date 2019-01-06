package yandexOAuth

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/khades/servbot/config"
	"github.com/khades/servbot/donationSource"
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/httpSession"
)

type Service struct {
	config                *config.Config
	donationSourceService *donationSource.Service
}

func (service *Service) login(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession) {
	code := r.URL.Query().Get("code")
	key, err := service.getToken(code)
	if err == nil {
		service.donationSourceService.SetYandexKey(s.UserID, key)
		json.NewEncoder(w).Encode(httpAPI.OptionResponse{"OK"})
		return
	}
	httpAPI.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
}

func (service *Service) getToken(code string) (string, error) {
	apiUrl := "https://money.yandex.ru/oauth/authorize"
	data := url.Values{}
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", "https://servbot.khades.org/yandex/oauth")
	data.Set("client_id", service.config.YandexClientID)
	data.Set("client_secret", service.config.YandexClientSecret)
	client := &http.Client{Timeout: 5 * time.Second}
	r, _ := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	resp, err := client.Do(r)
	if err != nil {
		return "", err
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	result := yandexResponse{}
	marshallError := json.NewDecoder(resp.Body).Decode(&result)

	dump, dumpErr := httputil.DumpResponse(resp, true)
	if dumpErr == nil {
		log.Printf("Repsonse is %q", dump)
	}
	if marshallError != nil {
		return "", marshallError
	}
	return result.AccessToken, nil
}
