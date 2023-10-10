package bot

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/alfredosa/GoDiscordBot/config"

	"github.com/bwmarrin/discordgo"
)

var BotId string
var goBot *discordgo.Session
var rule *discordgo.AutoModerationRule

type Bot struct {
	Session *discordgo.Session
	Rule    *discordgo.AutoModerationRule
	BotID   string
}

func NewBot() (*Bot, error) {
	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, err
	}

	session.Identify.Intents |= discordgo.IntentAutoModerationExecution
	session.Identify.Intents |= discordgo.IntentMessageContent
	session.Identify.Intents |= discordgo.IntentGuildMessages

	u, err := session.User("@me")
	if err != nil {
		return nil, err
	}

	return &Bot{
		Session: session,
		BotID:   u.ID,
	}, nil
}

func (b *Bot) CreateMessageTriggeredModRule(name string, keyword string, regex string) (string, error) {
	enabled := true
	rule, err := b.Session.AutoModerationRuleCreate(config.GuildId, &discordgo.AutoModerationRule{
		Name:        name,
		EventType:   discordgo.AutoModerationEventMessageSend,
		TriggerType: discordgo.AutoModerationEventTriggerKeyword,
		TriggerMetadata: &discordgo.AutoModerationTriggerMetadata{
			KeywordFilter: []string{keyword},
			RegexPatterns: []string{regex},
		},

		Enabled: &enabled,
		Actions: []discordgo.AutoModerationAction{
			{Type: discordgo.AutoModerationRuleActionSendAlertMessage, Metadata: &discordgo.AutoModerationActionMetadata{ChannelID: config.ChannelID}},
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		goBot.AutoModerationRuleDelete(config.GuildId, rule.ID)
		return "", err
	}

	fmt.Println("Rule created with ID: " + rule.ID)
	return rule.ID, nil
}

func Start() error {
	goBot, err := NewBot()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	ruleID, err := goBot.CreateMessageTriggeredModRule("Fuck rule", "*fuck*", "(f|d)uck")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	defer goBot.Session.AutoModerationRuleDelete(config.GuildId, ruleID)

	BotId = goBot.BotID
	fmt.Println("BotId: " + BotId)
	goBot.Session.AddHandler(messageHandler)
	fmt.Println("messageHandler added")
	goBot.Session.AddHandlerOnce(automodarationHandler)
	fmt.Println("automodarationHandler added")

	err = goBot.Session.Open()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Bot is running!")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan
	fmt.Println("Bot is shutting down...")
	return goBot.Session.Close()
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotId {
		return
	}

	// if m.content contains botid (Mentions) and "ping" then send "pong!"
	if m.Content == "<@"+BotId+"> !ping" || m.Content == "<@"+BotId+"> ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong!")
	}

	if strings.Contains(m.Content, "<@"+BotId+"> !google ") {
		_, _ = s.ChannelMessageSend(m.ChannelID, "https://www.google.com/search?q="+PrepareURLSearch(m.Content))
	}

	// send image from youtube
	if strings.Contains(m.Content, "<@"+BotId+"> !youtube ") {
		_, _ = s.ChannelMessageSend(m.ChannelID, "https://www.youtube.com/results?search_query="+PrepareURLSearch(m.Content))
		embed := &discordgo.MessageEmbed{
			Title:       "Embed Title",
			URL:         "https://github.com/bwmarrin/discordgo",
			Description: "Embed Description",
			Timestamp:   "2021-05-28",
			Color:       0x78141b,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Inline field 1 title",
					Value:  "value 1",
					Inline: true,
				},
				{
					Name:   "Inline field 2 title",
					Value:  "value 2",
					Inline: true,
				},
				{
					Name:   "Regular field title",
					Value:  "value 3",
					Inline: false,
				},
				{
					Name:   "Regular field 2 title",
					Value:  "value 4",
					Inline: false,
				},
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}

	if m.Content == config.BotPrefix+"ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong!")
	}
}

func PrepareURLSearch(content string) string {
	stringValueRaw := strings.Replace(content, "<@"+BotId+"> !google ", "", 1)
	stringwoYoutube := strings.Replace(stringValueRaw, "<@"+BotId+"> !youtube ", "", 1)
	stringPrepared := strings.Replace(stringwoYoutube, " ", "+", -1)
	return stringPrepared
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
