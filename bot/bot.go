package bot

import (
	"fmt"

	"github.com/alfredosa/GoDiscordBot/config"

	"github.com/bwmarrin/discordgo"
)

var BotId string
var goBot *discordgo.Session

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	goBot.Identify.Intents |= discordgo.IntentAutoModerationExecution
	goBot.Identify.Intents |= discordgo.IntentMessageContent

	// add rule

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
	}

	BotId = u.ID

	goBot.AddHandler(messageHandler)

	goBot.AddHandler(automodarationHandler)

	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running!")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotId {
		return
	}

	// if m.content contains botid (Mentions) and "ping" then send "pong!"
	if m.Content == "<@"+BotId+"> !ping" || m.Content == "<@"+BotId+"> ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong!")
	}

	if m.Content == config.BotPrefix+"ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong!")
	}
}

func automodarationHandler(s *discordgo.Session, e *discordgo.AutoModerationActionExecution) {
	fmt.Println("AutoModarationHandler called")

}
