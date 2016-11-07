package repos

import (
	"time"

	cache "github.com/patrickmn/go-cache"
)

var cacheObject = cache.New(15*time.Minute, 30*time.Second)
