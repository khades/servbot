package models

// HTTPSession defines nasfasf
type HTTPSession struct {
	Username  string `valid:"required"`
	UserID    string
	Key       string `valid:"required"`
	AvatarURL string
}
