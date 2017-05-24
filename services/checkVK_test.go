package services

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/khades/servbot/models"
)

func TestMobster(t *testing.T) {
	log.Println("Staring test")
	id := models.VkGroupInfo{GroupName: "mob5tervk"}
	ParseVK(&id)
}

func TestWTF(t *testing.T) {
	loc, _ := time.LoadLocation("Europe/Moscow")
	nowTime := time.Unix(0, 1495650294*1000000000).In(loc)
	fmt.Println(nowTime.Format("Jan _2 15:04 MSK"))
}
