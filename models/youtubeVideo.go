package models

import (
	"strconv"
	"time"

	duration "github.com/khades/iso8601duration"
	"github.com/sirupsen/logrus"
)

// YoutubeVideo describes youtube video information, returned by youtube
type YoutubeVideo struct {
	PageInfo YTPageInfo `json:"pageInfo"`
	Items    []YTItem   `json:"items"`
}

// YTPageInfo describes statistics of YoutubeVideo response
type YTPageInfo struct {
	TotalResults int `json:"totalResults"`
}

// YTItem describes one parsed youtube video information
type YTItem struct {
	Snippet        YTSnippet        `json:"snippet"`
	ContentDetails YTContentDetails `json:"contentDetails"`
	Statistics     YTStatistics     `json:"statistics"`
}

// YTSnippet describes parsed "snippet" value of youtube video
type YTSnippet struct {
	Title string `json:"title"`
}

// YTContentDetails  describes parsed "contentDetails" value of youtube video
type YTContentDetails struct {
	Duration string `json:"duration"`
}

// YTStatistics  describes parsed "statistics" value of youtube video
type YTStatistics struct {
	ViewCount string `json:"viewCount"`
	Likes     string `json:"likeCount"`
	Dislikes  string `json:"dislikeCount"`
}

// GetViewCount returns view count of one youtube video
func (ytStatistics YTStatistics) GetViewCount() int64 {
	logger := logrus.WithFields(logrus.Fields{
		"package": "models",
		"feature": "youtube",
		"action":  "GetViewCount"})
	value, error := strconv.ParseInt(ytStatistics.ViewCount, 10, 64)
	if error != nil {
		logger.Infof("ERROR: " + error.Error())
	}
	return value
}

// GetDuration returns duration of youtube video
func (ytContentDetails YTContentDetails) GetDuration() (*time.Duration, error) {
	return duration.ParseString(ytContentDetails.Duration)
}
