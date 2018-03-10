package httpclient

import (
	"io"
	"net/http"
	"time"
)

// TwitchV5 defines twitch v5 api oriented http client
func TwitchV5(clientID string, method string, urlStr string, body io.Reader) (*http.Response, error) {
	var timeout = 5 * time.Second
	var client = http.Client{Timeout: timeout}
	req, error := http.NewRequest(method, urlStr, body)
	req.Header.Add("Client-ID", clientID)
	req.Header.Add("Accept", "application/vnd.twitchtv.v5+json")
	if error != nil {
		return nil, error
	}
	return client.Do(req)
}

// Get defines simple GET http client with 5 seconds timeout
func Get(url string) (*http.Response, error) {
	var timeout = 5 * time.Second
	var client = http.Client{Timeout: timeout}
	return client.Get(url)
}
