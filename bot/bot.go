package bot

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/alfredosa/GoDiscordBot/config"

	"github.com/bwmarrin/discordgo"
)

var BotId string

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
	var rule *discordgo.AutoModerationRule
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
			{Type: discordgo.AutoModerationRuleActionTimeout, Metadata: &discordgo.AutoModerationActionMetadata{Duration: 60}},
			{Type: discordgo.AutoModerationRuleActionSendAlertMessage, Metadata: &discordgo.AutoModerationActionMetadata{ChannelID: config.ChannelID}},
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		b.Session.AutoModerationRuleDelete(config.GuildId, rule.ID)
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

	ruleID, err := goBot.CreateMessageTriggeredModRule("rust c++ rule", "i like c++, *assembly*, *c++*, I will rewritte my whole codebase in rust", "^(?:[0-9]{1,3}\\.){3}[0-9]{1,3}$, ")

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
		preparedURL := PrepareURLSearch(m.Content)
		currentTime := time.Now()
		// Format the current qtime as a string using a specific layout
		formattedTime := currentTime.Format("2006-01-02 15:04:05")
		embed := &discordgo.MessageEmbed{
			Title:       strings.Replace(m.Content, "<@"+BotId+"> !youtube ", " ", 1),
			URL:         "https://www.youtube.com/results?search_query=" + preparedURL,
			Description: "Youtube Search",
			Timestamp:   formattedTime,
			Color:       0x78141b,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   preparedURL,
					Value:  "The search you performed",
					Inline: true,
				},
			},
			Image: &discordgo.MessageEmbedImage{
				URL: "https://revistabyte.es/wp-content/uploads/2022/07/que-es-un-desarrollador-de-go-y-como-convertirse-en-uno-696x416.jpg.webp",
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

	if e.UserID == BotId {
		return
	}

	action := "unknown"
	switch e.Action.Type {
	case discordgo.AutoModerationRuleActionBlockMessage:
		action = "block message"
	case discordgo.AutoModerationRuleActionSendAlertMessage:
		action = "Alert message"
	case discordgo.AutoModerationRuleActionTimeout:
		action = "timeout"
	}

	s.ChannelMessageSend(e.ChannelID, "You just triggered the forbidden words.\n"+
		"Please don't do that :), you triggered ("+action+")\n"+
		"You will be **timed out** for a minute if you do.\n"+
		"Message that triggered the rule: "+e.Content+"\n")
}
