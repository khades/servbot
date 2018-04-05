package models

import (
	"log"
	"testing"
)

func TestSongRequests(t *testing.T) {
	videos := SongRequests{
		SongRequest{
			VideoID: "1",
			UserID:  "a",
			Order:   1},
		SongRequest{
			VideoID: "2",
			UserID:  "b",
			Order:   2},
		SongRequest{
			VideoID: "3",
			UserID:  "a",
			Order:   3},
		SongRequest{
			VideoID: "4",
			UserID:  "c",
			Order:   4}}
	userID := "a"

	newArray, _ := videos.PullUsersLastRequest(&userID)
	log.Printf("%+v", newArray)

}


func TestSongRequests3(t *testing.T) {
	videos := SongRequests{
		SongRequest{
			VideoID: "1",
			UserID:  "a",
			Order:   1},
		SongRequest{
			VideoID: "2",
			UserID:  "b",
			Order:   2},
		SongRequest{
			VideoID: "3",
			UserID:  "a",
			Order:   3},
		SongRequest{
			VideoID: "4",
			UserID:  "c",
			Order:   4}}

	videoID := "2"

	newArray2, _ := videos.PullOneRequest(&videoID)
	log.Printf("%+v", newArray2)

}


func TestSongRequests2(t *testing.T) {
	videos := SongRequests{
		SongRequest{
			VideoID: "1",
			UserID:  "a",
			Order:   1},
		SongRequest{
			VideoID: "2",
			UserID:  "b",
			Order:   2},
		SongRequest{
			VideoID: "3",
			UserID:  "a",
			Order:   3},
		SongRequest{
			VideoID: "4",
			UserID:  "c",
			Order:   4}}

	videoID := "2"


	sorted, _ := videos.BubbleVideoUp(&videoID)
	log.Printf("%+v", sorted)
}