package httpclient

import (
	"io"
	"net/http"
	"time"
)

func TwitchV5(clientId string, method string, urlStr string, body io.Reader) (*http.Response, error) {
	var timeout = 5 * time.Second
	var client = http.Client{Timeout: timeout}
	req, error := http.NewRequest(method, urlStr, body)
	req.Header.Add("Client-ID", clientId)
	req.Header.Add("Accept", "application/vnd.twitchtv.v5+json")
	if error != nil {
		return nil, error
	}
	return client.Do(req)
}

func Get(url string) (*http.Response, error) {
	var timeout = 5 * time.Second
	var client = http.Client{Timeout: timeout}
	return client.Get(url)
}
