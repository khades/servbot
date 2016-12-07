package httpbackend

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"

	"goji.io/pat"
)

func template(w http.ResponseWriter, r *http.Request, s *models.HTTPSession) {
	channel := pat.Param(r, "channel")
	template := pat.Param(r, "template")
	if channel == "" || template == "" {
		writeJSONError(w, "URL is not valid", http.StatusBadRequest)
		return
	}
	log.Println(channel)
	result, error := repos.GetChannelTemplateWithHistory(&channel, &template)
	if error != nil {
		writeJSONError(w, error.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(*result)
}
