package bot

import (
	"log"

	"github.com/khades/servbot/repos"
)

func modHandler(channel *string, mods *[]string) {
	log.Println("We got mods")
	log.Println(*mods)
	users, error := repos.GetUsersID(mods)
	if error != nil {
		return
	}
	values, error := repos.GetUsersID(&[]string{*channel})
	channelID := (*values)[*channel]
	if error != nil || channelID == "" {
		return
	}
	userIDs := []string{}
	for _, id := range *users {
		userIDs = append(userIDs, id)
	}
	repos.PushMods(&channelID, &userIDs)
}
