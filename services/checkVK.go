package services

import (
	"encoding/json"
	"errors"
	"fmt"
	//"log"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/httpclient"
	"github.com/khades/servbot/models"
	"github.com/khades/servbot/repos"
	"github.com/sirupsen/logrus"
)

type responseItem struct {
	ID       int    `json:"id"`
	Owner    int    `json:"owner_id"`
	Text     string `json:"text"`
	IsPinned int    `json:"is_pinned"`
	Date     int    `json:"date"`
}

type vkResponse struct {
	Response response `json:"response"`
}

type response struct {
	Items []responseItem `json:"items"`
}

func short(s string, i int) string {
	runes := []rune(s)
	if len(runes) > i {
		return string(runes[:i])
	}
	return s
}

// CheckVK checks all vk groups of all channels on that instance of bot
func CheckVK() {
	logger := logrus.WithFields(logrus.Fields{
		"package": "services",
		"feature": "vk",
		"action":  "CheckVK("})
	logger.Debugf("Checking VK")
	if repos.Config.VkClientKey == "" {
		logger.Debugf("VK key is not set")
		return
	}
	channels, error := repos.GetVKEnabledChannels()
	if error != nil {
		return
	}
	for _, channel := range channels {
		checkOne(&channel)
	}
}
func checkOne(channel *models.ChannelInfo) {
	logger := logrus.WithFields(logrus.Fields{
		"package": "services",
		"feature": "vk",
		"action":  "checkOne"})
	logger.Debug("Checking group " + channel.VkGroupInfo.GroupName)
	result, parseError := ParseVK(&channel.VkGroupInfo)
	if parseError != nil {
		logger.Debug("ParseError " + parseError.Error())
		return
	}
	if result.LastMessageID == channel.VkGroupInfo.LastMessageID {
		return
	}
	repos.PushVkGroupInfo(&channel.ChannelID, result)
	if result.NotifyOnChange == false {
		return
	}

	logger.Debug("Sending message to channel")
	bot.IrcClientInstance.SendPublic(&models.OutgoingMessage{
		Channel: channel.Channel,
		Body:    "[VK https://vk.com/" + channel.VkGroupInfo.GroupName + "] " + result.LastMessageBody + " " + result.LastMessageURL})

}

// ParseVK gets latest vk group post
func ParseVK(vkInputGroupInfo *models.VkGroupInfo) (*models.VkGroupInfo, error) {
	vkGroupInfo := models.VkGroupInfo{GroupName: vkInputGroupInfo.GroupName,
		NotifyOnChange: vkInputGroupInfo.NotifyOnChange}
	url := "https://api.vk.com/method/wall.get?domain=" + vkInputGroupInfo.GroupName + "&filter=owner&count=2&v=5.60"
	if strings.HasPrefix(vkInputGroupInfo.GroupName, "club") {
		url = "https://api.vk.com/method/wall.get?owner_id=-" + strings.Replace(vkInputGroupInfo.GroupName, "club", "", -1) + "&filter=owner&count=2&v=5.60"

	}
	resp, error := httpclient.Get(url + "&access_token=" + repos.Config.VkClientKey)
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
		vkPost.Text = short(vkPost.Text, 297) + "..."
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
