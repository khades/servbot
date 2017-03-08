package httpbackend

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"

	"goji.io/pat"
)

type autoMessageCreationResponse struct {
	ID bson.ObjectId
}
type autoMessageListResponse struct {
	Channel      string               `json:"channel"`
	AutoMessages []models.AutoMessage `json:"autoMessages"`
}

func autoMessageList(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	result, error := repos.GetAutoMessages(channelID)
	if error != nil {
		writeJSONError(w, error.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(autoMessageListResponse{*channelName, *result})

}

func autoMessageGet(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	id := pat.Param(r, "messageID")

	if id == "" {
		writeJSONError(w, "messageID or channel variable is not defined", http.StatusUnprocessableEntity)
		return
	}
	result, error := repos.GetAutoMessage(&id, channelID)
	if error != nil {
		writeJSONError(w, error.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func autoMessageCreate(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	decoder := json.NewDecoder(r.Body)
	var request models.AutoMessageUpdate
	err := decoder.Decode(&request)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "Invalid entry", http.StatusUnprocessableEntity)
		return
	}
	request.User = s.Username
	request.UserID = s.UserID

	request.ChannelID = *channelID

	id := repos.CreateAutoMessage(&request)

	json.NewEncoder(w).Encode(autoMessageCreationResponse{*id})

}

func autoMessageUpdate(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	id := pat.Param(r, "id")

	if id == "" || bson.IsObjectIdHex(id) == false {
		writeJSONError(w, "channel or id variable are not defined", http.StatusUnprocessableEntity)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var request models.AutoMessageUpdate
	err := decoder.Decode(&request)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "Invalid entry", http.StatusUnprocessableEntity)
		return
	}
	request.User = s.Username
	request.UserID = s.UserID
	request.ChannelID = *channelID
	request.ID = id
	repos.UpdateAutoMessage(&request)
}
