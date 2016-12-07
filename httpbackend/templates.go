package httpbackend

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"

	"goji.io/pat"
)

func templates(w http.ResponseWriter, r *http.Request, s *models.HTTPSession) {
	channel := pat.Param(r, "channel")
	if channel == "" {
		writeJSONError(w, "URL is not valid", http.StatusBadRequest)
		return
	}
	log.Println(channel)
	templates, error := repos.GetChannelTemplates(&channel)
	if error != nil {
		writeJSONError(w, error.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(*templates)
}
