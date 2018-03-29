package httpbackend

import (
	"encoding/json"
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type templatesResponse struct {
	Templates []models.TemplateInfo `json:"templates"`
	IsMod     bool                  `json:"isMod"`
}

func templates(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	channelInfo, error := repos.GetChannelInfo(channelID)
	if error != nil {
		writeJSONError(w, "That channel is not defined", http.StatusForbidden)
		return
	}
	templates, error := repos.GetChannelTemplates(channelID)
	if error != nil && error.Error() != "not found" {
		writeJSONError(w, error.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(templatesResponse{templates, channelInfo.GetIfUserIsMod(&s.UserID)})
}
