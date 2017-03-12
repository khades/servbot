package httpbackend

import (
	"net/http"
	"time"

	"github.com/JanBerktold/sse"
)

func timeTicker(w http.ResponseWriter, r *http.Request) {
	ticker := time.NewTicker(time.Second * 15)
	conn, _ := sse.Upgrade(w, r)
	for {
		tick := <-ticker.C
		conn.WriteJson(tick)
	}

}
