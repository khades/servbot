package vkGroupSchedule

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/config"
	"github.com/khades/servbot/twitchIRCClient"
	"github.com/khades/servbot/utils"
	"github.com/sirupsen/logrus"
)

type Service struct {
	config             *config.Config
	channelInfoService *channelInfo.Service
	twitchIRCClient    *twitchIRCClient.TwitchIRCClient
}

// Check checks all vk groups of all channels on that instance of bot
func (service *Service) Check() {
	logger := logrus.WithFields(logrus.Fields{
		"package": "services",
		"feature": "vk",
		"action":  "CheckVK"})

	logger.Debugf("Checking VK")
	if service.config.VkClientKey == "" {
		logger.Debugf("VK key is not set")
		return
	}
	channels, error := service.channelInfoService.GetVKEnabledChannels()
	if error != nil {
		return
	}
	for _, channel := range channels {
		service.checkOne(&channel)
	}
}

func (service *Service) checkOne(channel *channelInfo.ChannelInfo) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "services",
		"feature": "vk",
		"action":  "checkOne"})
	logger.Debug("Checking group " + channel.VkGroupInfo.GroupName)
	result, parseError := service.parseVK(&channel.VkGroupInfo)
	if parseError != nil {
		logger.Debug("ParseError " + parseError.Error())
		return
	}
	if result.LastMessageID == channel.VkGroupInfo.LastMessageID {
		return
	}
	service.channelInfoService.PushVkGroupInfo(&channel.ChannelID, result)
	if result.NotifyOnChange == false {
		return
	}

	logger.Debug("Sending message to channel")
	service.twitchIRCClient.SendPublic(&twitchIRCClient.OutgoingMessage{
		Channel: channel.Channel,
		Body:    "[VK https://vk.com/" + channel.VkGroupInfo.GroupName + "] " + result.LastMessageBody + " " + result.LastMessageURL})
}

// ParseVK gets latest vk group post
func (service *Service) parseVK(vkInputGroupInfo *channelInfo.VkGroupInfo) (*channelInfo.VkGroupInfo, error) {
	vkGroupInfo := channelInfo.VkGroupInfo{GroupName: vkInputGroupInfo.GroupName,
		NotifyOnChange: vkInputGroupInfo.NotifyOnChange}
	url := "https://api.vk.com/method/wall.get?domain=" + vkInputGroupInfo.GroupName + "&filter=owner&count=2&v=5.60"
	if strings.HasPrefix(vkInputGroupInfo.GroupName, "club") {
		url = "https://api.vk.com/method/wall.get?owner_id=-" + strings.Replace(vkInputGroupInfo.GroupName, "club", "", -1) + "&filter=owner&count=2&v=5.60"

	}
	var timeout = 5 * time.Second
	var client = http.Client{Timeout: timeout}
	resp, error := client.Get(url + "&access_token=" + service.config.VkClientKey)
	if error != nil {
		return &vkGroupInfo, error
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	vkResp := vkResponse{}
	marshallError := json.NewDecoder(resp.Body).Decode(&vkResp)
	if marshallError != nil {
		return &vkGroupInfo, marshallError
	}
	if len(vkResp.Response.Items) == 0 {
		return &vkGroupInfo, errors.New("not found")
	}
	vkPost := vkResp.Response.Items[0]
	if vkPost.IsPinned == 1 {
		if len(vkResp.Response.Items) == 1 {
			return &vkGroupInfo, errors.New("not found")
		}
		vkPost = vkResp.Response.Items[1]

	}

	vkPost.Text = strings.Replace(vkPost.Text, "\n", " ", -1)
	if utf8.RuneCountInString(vkPost.Text) > 300 {
		vkPost.Text = utils.Short(vkPost.Text, 297) + "..."
	}
	if vkPost.ID == 0 || vkPost.Owner == 0 {
		return nil, errors.New("VK Error")
	}
	vkGroupInfo.LastMessageID = vkPost.ID
	vkGroupInfo.LastMessageBody = vkPost.Text

	vkGroupInfo.LastMessageURL = fmt.Sprintf("https://vk.com/%s?w=wall%d_%d", vkInputGroupInfo.GroupName, vkPost.Owner, vkPost.ID)
	loc, _ := time.LoadLocation("Europe/Moscow")
	nowTime := time.Unix(0, int64(vkPost.Date)*1000000000).In(loc)
	vkGroupInfo.LastMessageDate = nowTime.Format("Jan _2 15:04 MSK")
	return &vkGroupInfo, nil
}
