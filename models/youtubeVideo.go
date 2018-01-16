package models

import (
	"log"
	"strconv"
	"time"

	duration "github.com/khades/iso8601duration"
)

type YoutubeVideo struct {
	PageInfo YTPageInfo `json:"pageInfo"`
	Items    []YTItem   `json:"items"`
}
type YTPageInfo struct {
	TotalResults int `json:"totalResults"`
}
type YTItem struct {
	Snippet        YTSnippet        `json:"snippet"`
	ContentDetails YTContentDetails `json:"contentDetails"`
	Statistics     YTStatistics     `json:"statistics"`
}
type YTSnippet struct {
	Title string `json:"title"`
}
type YTContentDetails struct {
	Duration string `json:"duration"`
}
type YTStatistics struct {
	ViewCount string `json:"viewCount"`
}

func (ytStatistics YTStatistics) GetViewCount() int64 {
	value, error := strconv.ParseInt(ytStatistics.ViewCount, 10, 64)
	if error != nil {
		log.Println("ERROR: " + error.Error())
	}
	return value
}
func (ytContentDetails YTContentDetails) GetDuration() (*time.Duration, error) {
	return duration.ParseString(ytContentDetails.Duration)
}
