package repos

func GetStreamersID() (*[]string, error) {
	users, error := GetUsersID(&Config.Channels)
	if error != nil {
		return nil, error
	}

	userIDs := []string{}
	for _, id := range *users {
		userIDs = append(userIDs, id)
	}
	return &userIDs, nil
}
