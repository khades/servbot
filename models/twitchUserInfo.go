package models

// TwitchUserInfo describes user information, returned by twitch
type TwitchUserInfo struct {
	ID           string `json:"id"`
	Login        string `json:"login"`
	DisplayName  string `json:"display_name"`
	ProfileImage string `json:"profile_image_url"`
}

