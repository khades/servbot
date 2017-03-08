package bot

import "github.com/khades/servbot/repos"

func modHandler(channel *string, mods *[]string) {
	users, error := repos.GetUsersID(mods)
	if error != nil {
		return
	}
	values, error := repos.GetUsersID(&[]string{*channel})
	channelID := (*values)[*channel]
	if error != nil {
		return
	}
	userIDs := []string{}
	for _, id := range *users {
		userIDs = append(userIDs, id)
	}
	repos.PushMods(&channelID, &userIDs)
}
