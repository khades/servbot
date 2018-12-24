package httpSession

import "time"

//db.log_events.createIndex( { "createdAt": 1 }, { expireAfterSeconds: 3600 } )
type httpSessionDBstruct struct {
	HTTPSession `bson:",inline"`
	CreatedAt          time.Time
}

type HTTPSession struct {
	Username  string `valid:"required"`
	UserID    string `valid:"required"`
	Key       string `valid:"required"`
	AvatarURL string
}
