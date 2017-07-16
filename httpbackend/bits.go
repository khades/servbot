package httpbackend

import (
	"encoding/json"
	"net/http"

	"goji.io/pat"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type bitsResponse struct {
	Bits    []models.UserBits `json:"bits"`
	Channel string            `json:"channel"`
}

func bits(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {

	bits, error := repos.GetBitsForChannel(channelID)
	if error != nil && error.Error() != "not found" {
		writeJSONError(w, error.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(bitsResponse{*bits, *channelName})
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
	json.NewEncoder(w).Encode(*bits)
}
