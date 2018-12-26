package youtubeAPIClient

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/khades/servbot/config"
)

type YouTubeAPIClient struct {
	config *config.Config
}

func (service *YouTubeAPIClient) get(id *string) (*http.Response, error) {
	url := "https://content.googleapis.com/youtube/v3/videos?id=" + *id + "&part=snippet%2CcontentDetails%2Cstatistics&key=" + service.config.YoutubeKey
	var timeout = 5 * time.Second
	var client = http.Client{Timeout: timeout}
	return client.Get(url)
}

func (service *YouTubeAPIClient) search(input *string) (*http.Response, error) {
	url := "https://content.googleapis.com/youtube/v3/search?type=video&q=" + url.QueryEscape(*input) + "&maxResults=1&part=snippet&key=" + service.config.YoutubeKey
	var timeout = 5 * time.Second
	var client = http.Client{Timeout: timeout}
	return client.Get(url)
}

func (service *YouTubeAPIClient) Get(id *string) (*YoutubeVideo, error) {
	if service.config.YoutubeKey == "" {
		return nil, errors.New("YT key is not set")
	}
	resp, error := service.get(id)
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

func (service *YouTubeAPIClient) Search(id *string) (*YoutubeVideo, error) {
	if service.config.YoutubeKey == "" {
		return nil, errors.New("YT key is not set")
	}
	resp, error := service.search(id)

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

	return service.Get(&ytVideo.Items[0].ID.VideoID)
}
