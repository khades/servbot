package httpAPI

import (
	"encoding/json"
	"net/http"

	"github.com/khades/servbot/httpSession"
	"gopkg.in/asaskevich/govalidator.v4"
)

type HTTPError struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}

func WriteJSONError(w http.ResponseWriter, message string, headerCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(headerCode)
	json.NewEncoder(w).Encode(&HTTPError{Code: headerCode, Message: message})
}

func validateSession(r *http.Request, httpSessionService *httpSession.Service) bool {
	cookie, cookieErr := r.Cookie("oauth")

	if cookieErr != nil || cookie == nil || cookie.Value == "" {
		return false
	}

	userInfo, userInfoError := httpSessionService.Get(&cookie.Value)
	if userInfoError != nil {
		return false

	}

	_, err := govalidator.ValidateStruct(userInfo)
	if err != nil {
		return false

	}
	return true
}

func Options(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(OptionResponse{"OK"})

}
