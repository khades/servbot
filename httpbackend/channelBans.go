package httpbackend

import (
	"encoding/json"
	"net/http"

	"github.com/khades/servbot/models"

	"github.com/khades/servbot/repos"

)

func channelBans(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	results,error := repos.GetChannelBans(channelID)
	if error != nil && error.Error() != "not found" {
		writeJSONError(w, error.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(&results)
}