package eventbus

import evbus "github.com/asaskevich/EventBus"

var EventBus = evbus.New()

func EventSub(channelID *string) string {
	return "sub:" + *channelID
}

func Subtrain(channelID *string) string {
	return "subtrain:" + *channelID
}

func Songrequest(channelID *string) string {
	return "songrequest:" + *channelID
}

var SubtrainBus = evbus.New()