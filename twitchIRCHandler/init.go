package twitchIRCHandler

import (
	"github.com/khades/servbot/autoMessage"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/channelLogs"
	"github.com/khades/servbot/subAlert"
	"github.com/khades/servbot/subday"
	"github.com/khades/servbot/subscriptionInfo"
	"github.com/khades/servbot/userResolve"
)

func Init(subdayService *subday.Service,
	channelInfoService *channelInfo.Service,
	subAlertService *subAlert.Service,
	channelLogsService *channelLogs.Service,
	autoMessageService *autoMessage.Service,
	userResolveService *userResolve.Service,
	subscriptionInfoService *subscriptionInfo.Service) *TwitchIRCHandler {
	return &TwitchIRCHandler{
		subdayService:      subdayService,
		channelInfoService: channelInfoService,
		subAlertService:    subAlertService,
		channelLogsService: channelLogsService,
		autoMessageService: autoMessageService,
		userResolveService: userResolveService,
		subscriptionInfoService: subscriptionInfoService,
	}
}
