package httpbackend

import (
	"encoding/json"
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


func template(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {

	commandName := pat.Param(r, "commandName")
	if commandName == "" {
		writeJSONError(w, "URL is not valid", http.StatusBadRequest)
		return
	}

	result, error := repos.GetChannelTemplateWithHistory(channelID, &commandName)
	if error != nil && error.Error() == "not found" {
		result.ShowOffline = true
		result.ShowOnline = true
	}

	result.ChannelID = *channelID
	result.CommandName = commandName

	json.NewEncoder(w).Encode(*result)
}

func putTemplate(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	decoder := json.NewDecoder(r.Body)
	var request models.TemplateInfoBody
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
	request.AliasTo = commandName

	templateError := repos.SetChannelTemplate(&s.Username, &s.UserID, channelID, &commandName, &request)

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

	repos.SetChannelTemplateAlias(&s.Username, &s.UserID, channelID, &commandName, &request.AliasTo)

	json.NewEncoder(w).Encode(optionResponse{"OK"})

}
