package services

import (
	"encoding/json"
	"html"
	"log"

	"github.com/khades/servbot/httpclient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type dubTrack struct {
	Message string
	Data    struct {
		SongInfo struct {
			Name string
		}
	}
}

// CheckDubTrack checks last playing track
func CheckDubTrack() {
	channels, error := repos.GetDubTrackEnabledChannels()
	if error != nil {
		return
	}
	for _, channel := range *channels {
		checkOneDubTrack(&channel)
	}
}

func checkOneDubTrack(channel *models.ChannelInfo) {
	status := models.DubTrack{ID: channel.DubTrack.ID}
	defer repos.PushDubTrack(&channel.ChannelID, &status)
	//log.Printf("Checking %s twitchDj track \n", channel.Channel)
	resp, error := httpclient.Get("https://api.dubtrack.fm/room/" + channel.DubTrack.ID + "/playlist/active")

	if error != nil {
		log.Println(error)
		return
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	track := dubTrack{}
	marshallError := json.NewDecoder(resp.Body).Decode(&track)
	if marshallError != nil {
		log.Println(marshallError)
		return
	}
	log.Println(track)
	if track.Message == "OK" && track.Data.SongInfo.Name != "" {
		status.Playing = true
		status.Track = html.UnescapeString(track.Data.SongInfo.Name)
	}
}
