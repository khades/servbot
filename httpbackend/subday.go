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

func subdayList(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	results, error := repos.GetSubdays(channelID)
	if error != nil {
		writeJSONError(w, error.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(&results)
}

// func subdayLast(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
// 	channelInfo, error := repos.GetChannelInfo(channelID)
// 	if error != nil {
// 		writeJSONError(w, "That channel is not defined", http.StatusForbidden)
// 		return
// 	}
// 	if channelInfo.GetIfUserIsMod(&s.UserID) == true {
// 		result, error := repos.GetLastSubdayMod(channelID)
// 		if error != nil {
// 			writeJSONError(w, error.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		object := subdayWithMod{result, true}
// 		json.NewEncoder(w).Encode(object)
// 	} else {
// 		result, error := repos.GetLastSubday(channelID)
// 		if error != nil {
// 			writeJSONError(w, error.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		object := subdayWithModNoWinners{result, false}

// 		json.NewEncoder(w).Encode(object)
// 	}

// }

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
		if (id != "last") {
			result, error = repos.GetSubdayByIdMod(&id)
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
		if (id != "last") {
			result, error = repos.GetSubdayById(&id)
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
func subdayPullWinner (w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
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
	repos.CloseSubday(channelID, &id)
	json.NewEncoder(w).Encode(optionResponse{"OK"})

}
