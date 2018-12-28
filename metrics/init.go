package metrics

import "github.com/prometheus/client_golang/prometheus"

func Init() *Service {
	service := &Service{
		chatMessageCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "servbot_irc_messages_per_channel",
			Help: "How many IRC messages sent, partitioned by channel name OR \"private\" is whispered.",
		}, []string{"channel"}),
		chatSchedulerCounter: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "servbot_irc_messages_throttled",
			Help: "How many IRC messages get throttled after send queue processed.",
		}),
		twitchAPICounter: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "servbot_twitch_api_total",
			Help: "How many Twitch API requests done",
		}),
		twitchAPIPerRequestCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "servbot_twitch_api_per_feature",
			Help: "How many Twitch API requests done, partitioned by resource.",
		}, []string{"resource"}),
		channelInfoRequests: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "servbot_channel_info_total",
			Help: "How many channelInfo.Get requests are called.",
		}),
		channelInfoPerChannelRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "servbot_channel_info_per_channel",
			Help: "How many channelInfo.Get requests are called, partitioned by channel name.",
		}, []string{"channel"}),
		httpAPIChannelCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "servbot_httpapi_channel_calls",
			Help: "How many times httpAPI called for channelInfo, partitioned by channel name.",
		}, []string{"channel"}),
		httpAPIUserCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "servbot_httpapi_user_calls",
			Help: "How many times httpAPI called for user info, partitioned by username.",
		}, []string{"user"}),
	}

	prometheus.MustRegister(service.chatMessageCounter)
	prometheus.MustRegister(service.chatSchedulerCounter)
	prometheus.MustRegister(service.twitchAPICounter)
	prometheus.MustRegister(service.twitchAPIPerRequestCounter)
	prometheus.MustRegister(service.channelInfoRequests)
	prometheus.MustRegister(service.channelInfoPerChannelRequests)
	prometheus.MustRegister(service.httpAPIChannelCounter)
	prometheus.MustRegister(service.httpAPIUserCounter)
	return service
}
