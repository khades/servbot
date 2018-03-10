package httpbackend

import (
	"encoding/json"
	"net/http"

	"goji.io/pat"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func bits(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	search := ""

	bits, error := repos.GetBitsForChannel(channelID, &search)
	if error != nil && error.Error() != "not found" {
		writeJSONError(w, error.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(bits)
}

func bitsSearch(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	search := pat.Param(r, "search")

	bits, error := repos.GetBitsForChannel(channelID, &search)
	if error != nil && error.Error() != "not found" {
		writeJSONError(w, error.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(bits)
}

func userbits(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	userID := pat.Param(r, "userID")

	if userID == "" {
		writeJSONError(w, "userID variable is not defined", http.StatusUnprocessableEntity)
		return
	}
	bits, error := repos.GetBitsForChannelUser(channelID, &userID)
	if error != nil && error.Error() != "not found" {
		writeJSONError(w, error.Error(), http.StatusNotFound)
		return
	}
	if error != nil && error.Error() == "not found" {

		userNameResult := ""

		bitsEmptyResult := models.UserBitsWithHistory{
			UserBits: models.UserBits{
				ChannelID: *channelID,
				UserID:    userID,
				User:      userNameResult},
			History: []models.UserBitsHistory{}}
		json.NewEncoder(w).Encode(bitsEmptyResult)
		return
	}

	json.NewEncoder(w).Encode(*bits)
}
