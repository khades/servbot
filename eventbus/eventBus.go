package eventbus

import observable "github.com/GianlucaGuarini/go-observable"

var EventBus = observable.New()

func EventSub(channelID *string) string {
	return "sub:" + *channelID
}

var SubtrainBus = observable.New()