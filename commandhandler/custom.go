package commandhandler

import (
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/chatMessage"
	"strings"
	"unicode/utf8"

	"html"

	"github.com/cbroglie/mustache"
	"github.com/khades/servbot/models"
)

// custom handler checks if input command has template and then fills it in with mustache templating and sends to a specified/user
func  (service *CommandHandler) custom(channelInfo *channelInfo.ChannelInfo, chatMessage *chatMessage.ChatMessage, chatCommand models.ChatCommand) {

	// Here's most of magic
	templateObject := &templateExtendedObject{ChannelInfo: *channelInfo, IsMod: chatMessage.IsMod, CommandBody: chatCommand.Body, CommandBodyIsEmpty: chatCommand.Body == ""}
	templateObject.IsMod = chatMessage.IsMod
	templateObject.IsSub = chatMessage.IsSub
	templateObject.UserID = chatMessage.UserID
	templateObject.User = chatMessage.User
	templateObject.twitchIRCClient = service.twitchIRCClient
	template, err := service.templateService.Get(&chatMessage.ChannelID, &chatCommand.Command)

	templateObject.subdayService = service.subdayService
	templateObject.followersService = service.followersService
	templateObject.songRequestService = service.songRequestService

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
		service.twitchIRCClient.SendPublic(&models.OutgoingMessage{
			Channel: chatMessage.Channel,
			User:    templateObject.User,
			Body:    html.UnescapeString(message)})
		return
	}

	service.twitchIRCClient.SendDebounced(models.OutgoingDebouncedMessage{
		Message: models.OutgoingMessage{
			Channel: chatMessage.Channel,
			User:    templateObject.User,
			Body:    html.UnescapeString(message)},
		Command:    template.AliasTo,
		RedirectTo: redirectTo})

}

