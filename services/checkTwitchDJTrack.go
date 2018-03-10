package services

import (
	"encoding/json"
	"html"

	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/httpclient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
)

type tdjTrack struct {
	Title string
}

// CheckTwitchDJTrack checks last playing track
func CheckTwitchDJTrack() {
	channels, error := repos.GetTwitchDJEnabledChannels()
	if error != nil {
		return
	}
	for _, channel := range channels {
		checkOneTwitchDJTrack(&channel)
	}
}

func checkOneTwitchDJTrack(channel *models.ChannelInfo) {

	status := models.TwitchDJ{ID: channel.TwitchDJ.ID}
	defer repos.PushTwitchDJ(&channel.ChannelID, &status)
	//log.Printf("Checking %s twitchDj track \n", channel.Channel)
	resp, error := httpclient.Get("https://twitch-dj.ru/includes/back.php?func=get_track&channel=" + channel.TwitchDJ.ID)

	if error != nil {
		//log.Println(error)
		return
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	track := tdjTrack{}
	marshallError := json.NewDecoder(resp.Body).Decode(&track)
	if marshallError != nil {
		//log.Println(marshallError)
		return
	}
	if track.Title != "" {
		status.Playing = true
		status.Track = html.UnescapeString(track.Title)
	}
	if status.Playing == false {
		return
	}
	if channel.TwitchDJ.NotifyOnChange == true {
		if status.Playing == true && channel.TwitchDJ.Track != status.Track {
			bot.IrcClientInstance.SendPublic(&models.OutgoingMessage{
				Channel: channel.Channel,
				Body:    "[TwitchDJ] Now Playing: " + status.Track})
		}
	}
}
