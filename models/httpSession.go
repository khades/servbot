package models

// HTTPSession defines nasfasf
type HTTPSession struct {
	Username string `valid:"required"`
	Key      string `valid:"required"`
}
