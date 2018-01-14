package eventbus

import evbus "github.com/asaskevich/EventBus"

var EventBus = evbus.New()

func EventSub(channelID *string) string {
	return "sub:" + *channelID
}

var SubtrainBus = evbus.New()