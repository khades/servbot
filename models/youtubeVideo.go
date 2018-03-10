package models

import (
	"log"
	"strconv"
	"time"

	duration "github.com/khades/iso8601duration"
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
}

// GetViewCount returns view count of one youtube video
func (ytStatistics YTStatistics) GetViewCount() int64 {
	value, error := strconv.ParseInt(ytStatistics.ViewCount, 10, 64)
	if error != nil {
		log.Println("ERROR: " + error.Error())
	}
	return value
}

// GetDuration returns duration of youtube video
func (ytContentDetails YTContentDetails) GetDuration() (*time.Duration, error) {
	return duration.ParseString(ytContentDetails.Duration)
}
