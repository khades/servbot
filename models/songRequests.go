package models

import (
	"sort"
	"time"
)

// SongRequest struct descibes one song request
type SongRequest struct {
	User    string        `json:"user"`
	UserID  string        `json:"userID"`
	Date    time.Time     `json:"date"`
	VideoID string        `json:"videoID"`
	Length  time.Duration `json:"length"`
	Title   string        `json:"title"`
	Order   int           `json:"order"`
}

// ChannelSongRequestSettings struct descibes current settings for songrequest on channel
type ChannelSongRequestSettings struct {
	OnlySubs           bool  `json:"onlySubs"`
	PlaylistLength     int   `json:"playlistLength"`
	MaxVideoLength     int   `json:"maxVideoLength"`
	MaxRequestsPerUser int   `json:"maxRequestsPerUser"`
	VideoViewLimit     int64 `json:"videoViewLimit"`
}

// SongRequests is type alias to array of songrequests
type SongRequests []SongRequest

func (requests SongRequests) Len() int {
	return len(requests)
}

func (requests SongRequests) Less(i, j int) bool {
	return requests[i].Order < requests[j].Order
}

func (requests SongRequests) Swap(i, j int) {
	requests[i], requests[j] = requests[j], requests[i]
}

// PullOneRequest removes one video from songrequests list and returns new list
func (requests SongRequests) PullOneRequest(videoID *string) (SongRequests, *SongRequest) {

	sort.Sort(requests)
	order := -1
	index := -1

	for itemIndex, request := range requests {
		if request.VideoID == *videoID {
			index = itemIndex
			order = request.Order
			break
		}
	}

	if index == -1 {
		return requests, nil

	}
	rejectedItem := requests[index]
	requests = append(requests[:index], requests[index+1:]...)
	for itemIndex := range requests {
		if (order < requests[itemIndex].Order) {
			requests[itemIndex].Order = requests[itemIndex].Order - 1

		}
	}

	return requests, &rejectedItem 
}

// PullUsersLastRequest removes last video from user and returns new list
func (requests SongRequests) PullUsersLastRequest(userID *string) (SongRequests, *SongRequest) {
	sort.Sort(sort.Reverse(requests))

	videoIndex := len(requests)
	for index, value := range requests {

		if value.UserID == *userID {

			videoIndex = index
			break
		}
	}

	if videoIndex < len(requests) {
		videoID := requests[videoIndex].VideoID

		return requests.PullOneRequest(&videoID)
	}

	return requests, nil
}

// BubbleVideoUp prioritises one video
func (requests SongRequests) BubbleVideoUp(videoID *string) (SongRequests, bool) {

	videoIndex := len(requests)

	for index, value := range requests {
		if value.VideoID == *videoID {
			videoIndex = index
			break
		}
	}

	if videoIndex >= len(requests) {

		return requests, false
	}

	for itemIndex := range requests {
		if requests[itemIndex].Order < requests[videoIndex].Order {
			requests[itemIndex].Order = requests[itemIndex].Order + 1
		}

	}

	requests[videoIndex].Order = 1

	return requests, true
}

// ChannelSongRequest describes all info about song request for channel
type ChannelSongRequest struct {
	ChannelID string                     `json:"channelID"`
	Settings  ChannelSongRequestSettings `json:"settings"`
	Requests  SongRequests               `json:"requests"`
}
