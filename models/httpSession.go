package models

// HTTPSession struct describes session information of user
type HTTPSession struct {
	Username  string `valid:"required"`
	UserID    string `valid:"required"`
	Key       string `valid:"required"`
	AvatarURL string
}
