package httpbackend

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"

	"goji.io/pat"
)

type templatePushRequest struct {
	Template string `json:"template"`
}
type aliasToRequest struct {
	AliasTo string `json:"aliasTo"`
}
type templateResponse struct {
	Template models.TemplateInfoWithHistory `json:"template"`
	Channel  string                         `json:"channel"`
}

func template(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	commandName := pat.Param(r, "commandName")
	if commandName == "" {
		writeJSONError(w, "URL is not valid", http.StatusBadRequest)
		return
	}
	log.Println(channel)
	result, _ := repos.GetChannelTemplateWithHistory(channelID, &commandName)

	json.NewEncoder(w).Encode(templateResponse{*result, *channelName})
}

func putTemplate(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	decoder := json.NewDecoder(r.Body)
	var request templatePushRequest
	err := decoder.Decode(&request)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	commandName := pat.Param(r, "commandName")
	if commandName == "" {
		writeJSONError(w, "URL is not valid", http.StatusBadRequest)
		return
	}
	templateError := repos.TemplateCache.UpdateTemplate(&s.Username, &s.UserID, channelID, &commandName, &request.Template)

	if templateError != nil {
		writeJSONError(w, templateError.Error(), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w).Encode(optionResponse{"OK"})

}

func aliasTemplate(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	decoder := json.NewDecoder(r.Body)
	var request aliasToRequest
	err := decoder.Decode(&request)

	if err != nil {
		writeJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	commandName := pat.Param(r, "commandName")

	if commandName == "" {
		writeJSONError(w, "URL is not valid", http.StatusBadRequest)
		return
	}

	repos.TemplateCache.SetAliasto(&s.Username, &s.UserID, channelID, &commandName, &request.AliasTo)

	json.NewEncoder(w).Encode(optionResponse{"OK"})

}
