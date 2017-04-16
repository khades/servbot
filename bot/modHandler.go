package bot

import "github.com/khades/servbot/repos"

func modHandler(channel *string, mods *[]string) {
	filteredMods := []string{}
	for _, mod := range *mods {
		if mod != "" {
			filteredMods = append(filteredMods, mod)
		}

	}
	users, error := repos.GetUsersID(&filteredMods)
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
