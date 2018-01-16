package httpclient
import (
	"time"
	"net/http"
	"log"
)
func YoutubeVideo (id *string, ytKey *string) (*http.Response, error)  {
	url:= "https://content.googleapis.com/youtube/v3/videos?id="+*id+"&part=snippet%2CcontentDetails%2Cstatistics&key="+*ytKey
	log.Println(url)
	var timeout = 5 * time.Second
	var client = http.Client{Timeout: timeout}
	return client.Get(url)
}