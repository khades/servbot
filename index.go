package main

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/khades/servbot/bot"
	"github.com/khades/servbot/repos"
)

//InsertA is some shittiest function
func main() {
	govalidator.SetFieldsRequiredByDefault(true)
	fmt.Println("bot username: ", repos.Config.BotUserName)

	//repos.Db.C("testing").Insert(&Person{"Ale", "+55 53 8116 9639"})
	bot.Start()
}
