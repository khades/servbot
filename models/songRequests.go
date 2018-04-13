package models

import (
	"sort"
	"time"
)

// SongRequest struct descibes one song request
type SongRequest struct {
	User     string        `json:"user"`
	UserID   string        `json:"userID"`
	Date     time.Time     `json:"date"`
	VideoID  string        `json:"videoID"`
	Length   time.Duration `json:"length"`
	Title    string        `json:"title"`
	Order    int           `json:"order"`
	Tags     []TagRecord   `json:"tags"`
	Views    int64         `json:"views"`
	Likes    int64         `json:"likes"`
	Dislikes int64         `json:"dislikes"`
}

// ChannelSongRequestSettings struct descibes current settings for songrequest on channel
type ChannelSongRequestSettings struct {
	OnlySubs           bool     `json:"onlySubs"`
	PlaylistLength     int      `json:"playlistLength"`
	MaxVideoLength     int      `json:"maxVideoLength"`
	MaxRequestsPerUser int      `json:"maxRequestsPerUser"`
	VideoViewLimit     int64    `json:"videoViewLimit"`
	MoreLikes          bool     `json:"moreLikes"`
	AllowOffline       bool     `json:"allowOffline"`
	Volume             int      `json:"volume"`
	BannedTags         []string `json:"bannedTags"`
	SkipIfTagged       bool     `json:"skipIfTagged"`
}

type TagRecord struct {
	Tag    string `json:"tag"`
	User   string `json:"user"`
	UserID string `json:"userID"`
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
		if order < requests[itemIndex].Order {
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
func (requests SongRequests) BubbleVideoUp(videoID *string, position int) (SongRequests, bool) {

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

	if requests[videoIndex].Order == position {
		return requests, false
	}

	for itemIndex := range requests {
		// moving order2 request to order4, order 3 and 4 should be upped, its orders are greater than original order, but equals of less than needed order
		if position >= requests[itemIndex].Order && requests[itemIndex].Order > requests[videoIndex].Order {
			requests[itemIndex].Order = requests[itemIndex].Order - 1
		}

		// moving order4 request to order2, order 2 and 3 should be downed, its orders are lesser than original order, but equals of greater than needed order
		if position <= requests[itemIndex].Order && requests[itemIndex].Order < requests[videoIndex].Order {
			requests[itemIndex].Order = requests[itemIndex].Order + 1
		}

	}

	requests[videoIndex].Order = position

	return requests, true
}

// ChannelSongRequest describes all info about song request for channel
type ChannelSongRequest struct {
	ChannelID string                     `json:"channelID"`
	Settings  ChannelSongRequestSettings `json:"settings"`
	Requests  SongRequests               `json:"requests"`
}

type TaggedVideoResult struct {
	User                     string
	Length                   time.Duration
	Title                    string
	ChannelID                string
	RemovedYoutubeRestricted bool
	RemovedTwitchRestricted  bool
	RemovedChannelRestricted bool
	RemovedTagRestricted     bool
	Tag                      string
}

type SongRequestAddResult struct {
	YoutubeRestricted bool
	TwitchRestricted  bool
	ChannelRestricted bool
	TagRestricted     bool
	Offline           bool
	PlaylistIsFull    bool
	AlreadyInPlaylist bool
	TooManyRequests   bool
	InvalidLink       bool
	NothingFound      bool
	InternalError     bool
	TooLong           bool
	TooLittleViews    bool
	MoreDislikes      bool
	Success           bool
	Title             string
	Length            time.Duration
	LengthStr         string
	Tag               string
}
