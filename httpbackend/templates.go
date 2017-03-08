package httpbackend

import (
	"encoding/json"
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)


type templatesResponse struct {
	Templates []models.TemplateInfo `json:"templates"`
	Channel   string                `json:"channel"`
}

func templates(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {

	templates, error := repos.GetChannelTemplates(channelID)
	if error != nil {
		writeJSONError(w, error.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(templatesResponse{*templates, *channelName})
}
