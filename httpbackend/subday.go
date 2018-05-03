package httpbackend

import (
	"encoding/json"
	"net/http"

	"goji.io/pat"

	"github.com/khades/servbot/models"

	"github.com/khades/servbot/repos"
)

type subdayWithMod struct {
	*models.Subday
	IsMod bool `json:"isMod"`
}
type subdayWithModNoWinners struct {
	*models.SubdayNoWinners
	IsMod bool `json:"isMod"`
}
type subdayCreateStruct struct {
	Name     string `json:"name"`
	SubsOnly bool   `json:"subsOnly"`
}
type subdayCreateResp struct {
	ID string `json:"id"`
}

func subdayCreate(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	decoder := json.NewDecoder(r.Body)
	var request subdayCreateStruct
	err := decoder.Decode(&request)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	created, id := repos.CreateNewSubday(channelID, request.SubsOnly, &request.Name)
	if created == false {
		writeJSONError(w, "Subday already exists", http.StatusNotAcceptable)
		return
	}
	json.NewEncoder(w).Encode(subdayCreateResp{id.Hex()})

}

func subdayList(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	results, error := repos.GetSubdays(channelID)
	if error != nil {
		writeJSONError(w, error.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(&results)
}

func subdayByID(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	channelInfo, error := repos.GetChannelInfo(channelID)
	if error != nil {
		writeJSONError(w, "That channel is not defined", http.StatusForbidden)
		return
	}
	id := pat.Param(r, "subdayID")
	if id == "" {
		writeJSONError(w, "subday id is not found", http.StatusNotFound)
		return
	}
	if channelInfo.GetIfUserIsMod(&s.UserID) == true {
		var result *models.Subday
		if id != "last" {
			result, error = repos.GetSubdayByIDMod(&id)
			if error != nil {
				writeJSONError(w, error.Error(), http.StatusInternalServerError)
				return
			}

		} else {
			result, error = repos.GetLastSubdayMod(channelID)
			if error != nil {
				writeJSONError(w, error.Error(), http.StatusInternalServerError)
				return
			}
		}

		object := subdayWithMod{result, true}

		json.NewEncoder(w).Encode(object)

	} else {
		var result *models.SubdayNoWinners
		if id != "last" {
			result, error = repos.GetSubdayByID(&id)
			if error != nil {
				writeJSONError(w, error.Error(), http.StatusInternalServerError)
				return
			}

		} else {
			result, error = repos.GetLastSubday(channelID)
			if error != nil {
				writeJSONError(w, error.Error(), http.StatusInternalServerError)
				return
			}
		}
		object := subdayWithModNoWinners{result, false}

		json.NewEncoder(w).Encode(object)
	}
}
func subdayRandomize(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	id := pat.Param(r, "subdayID")
	if id == "" {
		writeJSONError(w, "subday id is not found", http.StatusNotFound)
		return
	}
	winner := repos.PickRandomWinnerForSubday(channelID, &id)
	if winner != nil {
		json.NewEncoder(w).Encode(winner)
	} else {
		json.NewEncoder(w).Encode(optionResponse{"OK"})
	}

}
func subdayPullWinner(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	id := pat.Param(r, "subdayID")
	if id == "" {
		writeJSONError(w, "subday id is not found", http.StatusNotFound)
		return
	}
	user := pat.Param(r, "user")
	if user == "" {
		writeJSONError(w, "user", http.StatusNotFound)
		return
	}
	repos.SubdayPullWinner(&id, &user)
	json.NewEncoder(w).Encode(optionResponse{"OK"})

}
func subdayClose(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	id := pat.Param(r, "subdayID")
	if id == "" {
		writeJSONError(w, "subday id is not found", http.StatusNotFound)
		return
	}
	repos.CloseSubday(channelID, &id, &s.Username, &s.UserID)
	json.NewEncoder(w).Encode(optionResponse{"OK"})

}
