package subtrainAPI

import (
	"encoding/json"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/httpAPI"
	"github.com/khades/servbot/httpSession"
	"net/http"
)

type Service struct {
	channelInfoService *channelInfo.Service
}
func (service *Service) get(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfo *channelInfo.ChannelInfo) {
	json.NewEncoder(w).Encode(&channelInfo.SubTrain)
}

//func subtrainWidget(w http.ResponseWriter, r *http.Request, channelID *string) {
//	channelInfo, error := repos.Get(channelID)
//	if error != nil {
//		writeJSONError(w, error.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	json.NewEncoder(w).Encode(&channelInfo.SubTrain)
//}

//func subtrainWidgetEvents(w http.ResponseWriter, r *http.Request, channelID *string) {
//	websocketEventbusWriter(w, r, eventbus.Subtrain(channelID))
//}

func (service *Service) set(w http.ResponseWriter, r *http.Request, s *httpSession.HTTPSession, channelInfoStruct *channelInfo.ChannelInfo) {
	decoder := json.NewDecoder(r.Body)
	var request channelInfo.SubTrain
	err := decoder.Decode(&request)
	if err != nil {
		httpAPI.WriteJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	request.CurrentStreak = 0
	request.NotificationShown = false
	request.Users = []string{}
	service.channelInfoService.PutChannelSubtrainWeb(&channelInfoStruct.ChannelID, &request)
	json.NewEncoder(w).Encode(httpAPI.OptionResponse{"OK"})
}

//func subtrainEvents(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
//	websocketEventbusWriter(w, r, eventbus.Subtrain(channelID))
//}
