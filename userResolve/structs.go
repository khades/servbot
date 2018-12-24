package userResolve

import "time"

type usernameCache struct {
	UserID    string
	User      string
	CreatedAt time.Time
}
