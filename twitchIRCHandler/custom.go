package twitchIRCHandler

import (
	"github.com/cbroglie/mustache"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/chatMessage"
	"github.com/khades/servbot/twitchIRC"
	"github.com/khades/servbot/utils"
	"html"
	"strings"
	"unicode/utf8"
)

// custom handler checks if input command has template and then fills it in with mustache templating and sends to a specified/user
func (service *TwitchIRCHandler) custom(
	channelInfo *channelInfo.ChannelInfo,
	chatMessage *chatMessage.ChatMessage,
	chatCommand chatMessage.ChatCommand,
	twitchIRCClient *twitchIRC.Client) {

	// Here's most of magic
	templateObject := &extendedChannelInfo{
		ChannelInfo: *channelInfo,
		UserID:             chatMessage.UserID,
		User:               chatMessage.User,
		IsMod:              chatMessage.IsMod,
		IsSub:              chatMessage.IsSub,
		CommandBody:        chatCommand.Body,
		CommandBodyIsEmpty: chatCommand.Body == "",
		twitchIRCClient:    twitchIRCClient,
		followersService:   service.followersService,
		subdayService:      service.subdayService,
		songRequestService: service.songRequestService}

	template, err := service.templateService.Get(&chatMessage.ChannelID, &chatCommand.Command)

	if err != nil || template.Template == "" {
		return
	}

	message, renderError := mustache.Render(template.Template, templateObject)
	if renderError != nil {
		return
	}
	message = strings.TrimSpace(message)
	if utf8.RuneCountInString(message) > 400 {
		message = utils.Short(message, 397) + "..."
	}
	if message == "" {
		return
	}
	redirectTo := chatMessage.User

	//if chatCommand.Body != "" && !(template.StringRandomizer.Enabled == true && len(template.StringRandomizer.Strings) == 0) && template.PreventRedirect == false {
	if templateObject.PreventRedirect == false && chatCommand.Body != "" {
		if strings.HasPrefix(chatCommand.Body, "@") {
			redirectTo = chatCommand.Body[1:]
		} else {
			redirectTo = chatCommand.Body

		}
	}

	redirectTo = strings.Replace(redirectTo, "@", " @", -1)
	//if template.PreventDebounce == true {
	if templateObject.PreventDebounce == true {
		twitchIRCClient.SendPublic(&twitchIRC.OutgoingMessage{
			Channel: chatMessage.Channel,
			User:    templateObject.User,
			Body:    html.UnescapeString(message)})
		return
	}

	twitchIRCClient.SendDebounced(twitchIRC.OutgoingDebouncedMessage{
		Message: twitchIRC.OutgoingMessage{
			Channel: chatMessage.Channel,
			User:    templateObject.User,
			Body:    html.UnescapeString(message)},
		Command:    template.AliasTo,
		RedirectTo: redirectTo})

}
