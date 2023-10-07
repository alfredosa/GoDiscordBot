package main

import (
	"fmt"

	"github.com/alfredosa/GoDiscordBot/config"

	"github.com/alfredosa/GoDiscordBot/bot"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err)
		return
	}

	bot.Start()

	<-make(chan struct{})
	return
}
