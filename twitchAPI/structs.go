package twitchAPI

import "time"

type twitchUserInfo struct {
	ID          string `json:"_id"`
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
}

type usersLoginStruct struct {
	Users []twitchUserInfo `json:"users"`
}

type twitchFollower struct {
	UserID    string    `json:"from_id"`
	ChannelID string    `json:"to_id"`
	Date      time.Time `json:"followed_at"`
}

type twitchFollowerResponsePagination struct {
	Cursor string `json:"cursor"`
}

type twitchFollowers []twitchFollower

func (followers twitchFollowers) Len() int {
	return len(followers)
}

func (followers twitchFollowers) Less(i, j int) bool {
	return followers[i].Date.Before(followers[j].Date)
}

func (followers twitchFollowers) Swap(i, j int) {
	followers[i], followers[j] = followers[j], followers[i]
}

type twitchFollowerResponse struct {
	Total      int64                            `json:"total"`
	Pagination twitchFollowerResponsePagination `json:"pagination"`
	Followers  twitchFollowers                  `json:"data"`
}

type twitchGameStruct struct {
	GameID string `json:"id"`
	Game   string `json:"name"`
}

type twitchGameStructResponse struct {
	Data []twitchGameStruct `json:"data"`
}

// TwitchStreamStatus describes stream status returned by twitch api
type TwitchStreamStatus struct {
	GameID      string    `json:"game_id"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"started_at"`
	Title       string    `json:"title"`
	ViewerCount int       `json:"viewer_count"`
}

// TwitchUserInfo describes user information, returned by twitch
type TwitchUserInfo struct {
	ID           string `json:"id"`
	Login        string `json:"login"`
	DisplayName  string `json:"display_name"`
	ProfileImage string `json:"profile_image_url"`
}

type twitchStreamResponse struct {
	Data []TwitchStreamStatus `json:"data"`
}

type twitchUserRepsonse struct {
	Data []TwitchUserInfo `json:"data"`
}

type hub struct {
	Mode         string `json:"hub.mode"`
	Topic        string `json:"hub.topic"`
	Callback     string `json:"hub.callback"`
	LeaseSeconds string `json:"hub.lease_seconds"`
	Secret       string `json:"hub.secret"`
}

type twitchAPIKeyResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int64 `json:"expires_in"`
}

type ApiKeyResponse struct {
	AccessToken string
	ExpiresIn time.Time
}