package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Service struct {
	chatMessageCounter *prometheus.CounterVec
	chatSchedulerCounter prometheus.Counter
	twitchAPICounter prometheus.Counter
	twitchAPIPerRequestCounter  *prometheus.CounterVec
	channelInfoRequests prometheus.Counter
	channelInfoPerChannelRequests *prometheus.CounterVec
	httpAPIChannelCounter *prometheus.CounterVec
	httpAPIUserCounter *prometheus.CounterVec
}

func (service *Service) LogIRCMessageOnChannel (channel string) {
	service.chatMessageCounter.WithLabelValues(channel).Inc()
}

func (service *Service) LogIRCSchedulerPayload (count float64) {
	service.chatSchedulerCounter.Add(count)
}

func (service *Service) LogTwitchRequest() {
	service.twitchAPICounter.Inc()
}

func (service *Service) LogTwitchSpecificRequest(request string) {
	service.twitchAPIPerRequestCounter.WithLabelValues(request).Inc()
}

func (service *Service) LogChannelInfoRetrieval() {
	service.channelInfoRequests.Inc()
}

func (service *Service) LogChannelInfoRetrievalPerChannel(channel string) {
	service.channelInfoPerChannelRequests.WithLabelValues(channel).Inc()
}

func (service *Service) LogHTTPApiUserRequest(user string) {
	service.httpAPIUserCounter.WithLabelValues(user).Inc()
}

func (service *Service) LogHTTPApiChannelRequest(channel string) {
	service.httpAPIChannelCounter.WithLabelValues(channel).Inc()
}