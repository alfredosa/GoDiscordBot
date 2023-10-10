package bot

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alfredosa/GoDiscordBot/config"

	"github.com/bwmarrin/discordgo"
)

var BotId string
var goBot *discordgo.Session
var rule *discordgo.AutoModerationRule

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	goBot.Identify.Intents |= discordgo.IntentAutoModerationExecution
	goBot.Identify.Intents |= discordgo.IntentMessageContent

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
	}

	enabled := true
	rule, err := goBot.AutoModerationRuleCreate(config.GuildId, &discordgo.AutoModerationRule{
		Name:        "Auto Moderation example",
		EventType:   discordgo.AutoModerationEventMessageSend,
		TriggerType: discordgo.AutoModerationEventTriggerKeyword,
		TriggerMetadata: &discordgo.AutoModerationTriggerMetadata{
			KeywordFilter: []string{"*cat*"},
			RegexPatterns: []string{"(c|b)at"},
		},

		Enabled: &enabled,
		Actions: []discordgo.AutoModerationAction{
			{Type: discordgo.AutoModerationRuleActionBlockMessage},
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		goBot.AutoModerationRuleDelete(config.GuildId, rule.ID)
		panic(err)
	}

	fmt.Println("Rule created with ID: " + rule.ID)
	defer goBot.AutoModerationRuleDelete(config.GuildId, rule.ID)

	BotId = u.ID
	fmt.Println("BotId: " + BotId)
	goBot.AddHandler(messageHandler)
	fmt.Println("messageHandler added")
	goBot.AddHandlerOnce(automodarationHandler)
	fmt.Println("automodarationHandler added")

	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running!")

	defer goBot.Close()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	<-signalChannel

	fmt.Println("Received termination signal. Cleaning up...")
	os.Exit(0)
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
	s.ChannelMessageSend(e.ChannelID, "Congratulations! You have just triggered an auto moderation rule.\n"+
		"The current trigger can match anywhere in the word, so even if you write the trigger word as a part of another word, it will still match.\n"+
		"The rule has now been changed, now the trigger matches only in the full words.\n"+
		"Additionally, when you send a message, an alert will be sent to this channel and you will be **timed out** for a minute.\n")

	var counter int
	var counterMutex sync.Mutex
	goBot.AddHandler(func(s *discordgo.Session, e *discordgo.AutoModerationActionExecution) {
		action := "unknown"
		switch e.Action.Type {
		case discordgo.AutoModerationRuleActionBlockMessage:
			action = "block message"
		case discordgo.AutoModerationRuleActionSendAlertMessage:
			action = "send alert message into <#" + e.Action.Metadata.ChannelID + ">"
		case discordgo.AutoModerationRuleActionTimeout:
			action = "timeout"
		}

		counterMutex.Lock()
		counter++
		if counter == 1 {
			counterMutex.Unlock()
			s.ChannelMessageSend(e.ChannelID, "Nothing has changed, right? "+
				"Well, since separate gateway events are fired per each action (current is "+action+"), "+
				"you'll see a second message about an action pop up soon")
		} else if counter == 2 {
			counterMutex.Unlock()
			s.ChannelMessageSend(e.ChannelID, "Now the second ("+action+") action got executed.")
			s.ChannelMessageSend(e.ChannelID, "And... you've made it! That's the end of the example.\n"+
				"For more information about the automod and how to use it, "+
				"you can visit the official Discord docs: https://discord.dev/resources/auto-moderation or ask in our server: https://discord.gg/6dzbuDpSWY",
			)
		}
	})
}
