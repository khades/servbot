package twitchIRCHandler

import (
	"github.com/asaskevich/EventBus"
	"github.com/khades/servbot/autoMessage"
	"github.com/khades/servbot/channelInfo"
	"github.com/khades/servbot/channelLogs"
	"github.com/khades/servbot/followers"
	"github.com/khades/servbot/pubsub"
	"github.com/khades/servbot/songRequest"
	"github.com/khades/servbot/subAlert"
	"github.com/khades/servbot/subday"
	"github.com/khades/servbot/subscriptionInfo"
	"github.com/khades/servbot/template"
	"github.com/khades/servbot/userResolve"
)

func Init(subdayService *subday.Service,
	channelInfoService *channelInfo.Service,
	subAlertService *subAlert.Service,
	channelLogsService *channelLogs.Service,
	autoMessageService *autoMessage.Service,
	userResolveService *userResolve.Service,
	subscriptionInfoService *subscriptionInfo.Service,
	templateService *template.Service,
	followersService *followers.Service,
	songRequestService *songRequest.Service,
	eventBus EventBus.Bus,
	pubsub *pubsub.Client,
	// eventService *event.Service,
	// balanceService *balance.Service
) *TwitchIRCHandler {
	return &TwitchIRCHandler{
		subdayService:           subdayService,
		channelInfoService:      channelInfoService,
		subAlertService:         subAlertService,
		channelLogsService:      channelLogsService,
		autoMessageService:      autoMessageService,
		userResolveService:      userResolveService,
		subscriptionInfoService: subscriptionInfoService,
		templateService:         templateService,
		followersService:        followersService,
		songRequestService:      songRequestService,
		eventBus:                eventBus,
		pubsub:                  pubsub,
		// eventService:            eventService,
		// balanceService:          balanceService,
	}
}
