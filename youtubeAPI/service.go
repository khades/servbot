package youtubeAPI

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/khades/servbot/config"
)

type Client struct {
	config *config.Config
}

func (client *Client) get(id *string) (*http.Response, error) {
	url := "https://content.googleapis.com/youtube/v3/videos?id=" + *id + "&part=snippet%2CcontentDetails%2Cstatistics&key=" + client.config.YoutubeKey
	var timeout = 5 * time.Second
	var httpClient = http.Client{Timeout: timeout}
	return httpClient.Get(url)
}

func (client *Client) search(input *string) (*http.Response, error) {
	url := "https://content.googleapis.com/youtube/v3/search?type=video&q=" + url.QueryEscape(*input) + "&maxResults=1&part=snippet&key=" + client.config.YoutubeKey
	var timeout = 5 * time.Second
	var httpClient = http.Client{Timeout: timeout}
	return httpClient.Get(url)
}

func (client *Client) Get(id *string) (*YoutubeVideo, error) {
	if client.config.YoutubeKey == "" {
		return nil, errors.New("YT key is not set")
	}
	resp, error := client.get(id)
	if error != nil {
		return nil, error
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	var ytVideo YoutubeVideo

	marshallError := json.NewDecoder(resp.Body).Decode(&ytVideo)
	if marshallError != nil {
		return nil, marshallError
	}
	return &ytVideo, nil
}

func (client *Client) Search(id *string) (*YoutubeVideo, error) {
	if client.config.YoutubeKey == "" {
		return nil, errors.New("YT key is not set")
	}
	resp, error := client.search(id)

	if error != nil {
		return nil, error
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	var ytVideo YoutubeStringVideo

	marshallError := json.NewDecoder(resp.Body).Decode(&ytVideo)
	if marshallError != nil {
		return nil, marshallError
	}
	if ytVideo.PageInfo.TotalResults == 0 {
		return nil, errors.New("not found")
	}

	return client.Get(&ytVideo.Items[0].ID.VideoID)
}
