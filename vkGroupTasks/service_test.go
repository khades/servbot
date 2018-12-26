package vkGroupTasks

import (
	"fmt"
	"github.com/khades/servbot/channelInfo"
	"testing"
	"time"

)

func TestMobster(t *testing.T) {
	id := channelInfo.VkGroupInfo{GroupName: "mob5tervk"}
	Parse(&id)
}

func TestWTF(t *testing.T) {
	loc, _ := time.LoadLocation("Europe/Moscow")
	nowTime := time.Unix(0, 1495650294*1000000000).In(loc)
	fmt.Println(nowTime.Format("Jan _2 15:04 MSK"))
}
