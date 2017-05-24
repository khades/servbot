package httpbackend

import (
	"encoding/json"
	"net/http"

	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

func putVK(w http.ResponseWriter, r *http.Request, s *models.HTTPSession, channelID *string, channelName *string) {
	decoder := json.NewDecoder(r.Body)
	var request models.VkGroupInfo
	err := decoder.Decode(&request)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	vkUpdateRequest := models.VkGroupInfo{GroupName: request.GroupName, NotifyOnChange: request.NotifyOnChange}
	repos.PushVkGroupInfo(channelID, &vkUpdateRequest)
	json.NewEncoder(w).Encode(optionResponse{"OK"})
}
