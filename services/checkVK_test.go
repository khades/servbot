package services

import (
	"log"
	"testing"

	"github.com/khades/servbot/models"
)

func TestMobster(t *testing.T) {
	log.Println("Staring test")
	id := models.VkGroupInfo{GroupName: "mob5tervk"}
	ParseVK(&id)
}
