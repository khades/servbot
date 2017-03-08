package httpbackend

import (
	"encoding/json"
	"net/http"

	"github.com/khades/servbot/models"

	"github.com/khades/servbot/repos"

	"goji.io/pat"
)

func logsByUsername(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	user := pat.Param(r, "user")
	if user == "" {
		writeJSONError(w, "Bad request, no user or channel set", http.StatusBadRequest)
		return
	}
	userLogs, error := repos.GetUserMessageHistoryByUsername(&user, channelID)
	if error != nil {
		writeJSONError(w, error.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(*userLogs)
}

func logsByUserID(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	user := pat.Param(r, "userID")
	if user == "" {
		writeJSONError(w, "Bad request, no user or channel set", http.StatusBadRequest)
		return
	}
	userLogs, error := repos.GetUserMessageHistoryByUserID(&user, channelID)
	if error != nil {
		writeJSONError(w, error.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(*userLogs)
}
