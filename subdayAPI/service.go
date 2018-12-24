package subdayAPI

import (
	"encoding/json"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/httpSession"
	"github.com/khades/servbot/subday"
	"net/http"

	"goji.io/pat"
)

type Service struct {
	subdayService *subday.Service
}

func (service *Service) create(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	decoder := json.NewDecoder(r.Body)
	var request subdayCreateStruct
	err := decoder.Decode(&request)
	if err != nil {
		httpAPI.WriteJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	created, id := service.subdayService.Create(&channelInfo.ChannelID, request.SubsOnly, &request.Name)
	if created == false {
		httpAPI.WriteJSONError(w, "Subday already exists", http.StatusNotAcceptable)
		return
	}
	json.NewEncoder(w).Encode(subdayCreateResp{id.Hex()})

}

func (service *Service) list(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	results, error := service.subdayService.List(&channelInfo.ChannelID)
	if error != nil {
		httpAPI.WriteJSONError(w, error.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(&results)
}

func (service *Service) get(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	var error error;
	id := pat.Param(r, "subdayID")
	if id == "" {
		httpAPI.WriteJSONError(w, "subday id is not found", http.StatusNotFound)
		return
	}
	if channelInfo.GetIfUserIsMod(&s.UserID) == true {
		var result *subday.Subday
		if id != "last" {
			result, error = service.subdayService.GetForMod(&id)
			if error != nil {
				httpAPI.WriteJSONError(w, error.Error(), http.StatusInternalServerError)
				return
			}

		} else {
			result, error = service.subdayService.GetLastForMod(&channelInfo.ChannelID)
			if error != nil {
				httpAPI.WriteJSONError(w, error.Error(), http.StatusInternalServerError)
				return
			}
		}

		object := subdayWithMod{result, true}

		json.NewEncoder(w).Encode(object)

	} else {
		var result *subday.SubdayNoWinners
		if id != "last" {
			result, error = service.subdayService.Get(&id)
			if error != nil {
				httpAPI.WriteJSONError(w, error.Error(), http.StatusInternalServerError)
				return
			}

		} else {
			result, error = service.subdayService.GetLast(&channelInfo.ChannelID)
			if error != nil {
				httpAPI.WriteJSONError(w, error.Error(), http.StatusInternalServerError)
				return
			}
		}
		object := subdayWithModNoWinners{result, false}

		json.NewEncoder(w).Encode(object)
	}
}

func (service *Service) randomize(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	id := pat.Param(r, "subdayID")
	if id == "" {
		httpAPI.WriteJSONError(w, "subday id is not found", http.StatusNotFound)
		return
	}
	winner := service.subdayService.PickRandomWinner(&channelInfo.ChannelID, &id, false, false)
	if winner != nil {
		json.NewEncoder(w).Encode(winner)
	} else {
		json.NewEncoder(w).Encode(httpAPI.OptionResponse{"OK"})
	}
}

func (service *Service) randomizeSubs(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	id := pat.Param(r, "subdayID")
	if id == "" {
		httpAPI.WriteJSONError(w, "subday id is not found", http.StatusNotFound)
		return
	}
	winner := service.subdayService.PickRandomWinner(&channelInfo.ChannelID, &id, true, false)
	if winner != nil {
		json.NewEncoder(w).Encode(winner)
	} else {
		json.NewEncoder(w).Encode(httpAPI.OptionResponse{"OK"})
	}
}

func (service *Service) randomizeNonSubs(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	id := pat.Param(r, "subdayID")
	if id == "" {
		httpAPI.WriteJSONError(w, "subday id is not found", http.StatusNotFound)
		return
	}
	winner := service.subdayService.PickRandomWinner(&channelInfo.ChannelID, &id, false, true)
	if winner != nil {
		json.NewEncoder(w).Encode(winner)
	} else {
		json.NewEncoder(w).Encode(httpAPI.OptionResponse{"OK"})
	}
}

func (service *Service) pullWinner(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	id := pat.Param(r, "subdayID")
	if id == "" {
		httpAPI.WriteJSONError(w, "subday id is not found", http.StatusNotFound)
		return
	}
	user := pat.Param(r, "user")
	if user == "" {
		httpAPI.WriteJSONError(w, "user", http.StatusNotFound)
		return
	}
	service.subdayService.PullWinner(&id, &user)
	json.NewEncoder(w).Encode(httpAPI.OptionResponse{"OK"})

}

func (service *Service) close(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	id := pat.Param(r, "subdayID")
	if id == "" {
		httpAPI.WriteJSONError(w, "subday id is not found", http.StatusNotFound)
		return
	}
	service.subdayService.Close(&channelInfo.ChannelID, &id, &s.Username, &s.UserID)
	json.NewEncoder(w).Encode(httpAPI.OptionResponse{"OK"})

}
