package httpbackend

import (
	"encoding/json"
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func channelInfo(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	channelInfo, _ := repos.GetChannelInfo(channelID)
	// if error != nil {
	// 	writeJSONError(w, error.Error(), http.StatusInternalServerError)
	// 	return
	// }
	json.NewEncoder(w).Encode(&channelInfo)
}
