# Discord bot with Golang

In this guide we will create a simple Discord bot using Golang.

## 1. Create a Dicord Server 

It's important that you have control over the server, so you can add the bot to it. This is vital because the bot will need to be added to the server in order to work. You can't just add a bot to any server.

- You have to be an admin of the server. 

## 2. Create a Discord Application

- Go to https://discord.com/developers/applications
- Create a new application
- Create a new bot
- Configure the Bot (And create a new Token)
- Add a cool image for your bot
- get the Application ID

## 3. Add bot to your server

- Grab the details above (Client_id in your *Generatl Information*)
- go to https://discord.com/oauth2/authorize?client_id=[CLIENT_ID]&scope=bot
- select your server

## 4. Ready to get coding

create the following structure

```bash
├── config.json
├── main.go
├── bot
│   ├── bot.go
├── config
│   ├── config.go
```

Initialize your project, and install the following packages

```bash
go mod init github.com/[YOUR_USERNAME]/[YOUR_PROJECT_NAME]
go get github.com/bwmarrin/discordgo
```

## 5. json config

From step 1 and 2 you should have the following information

```json
{
    "token": "some_token_here",
    "BotPrefix": "!"
}
```

The bot prefix helps identify when a user is talking to the bot or not. 

## 5b. config/config.go

```go
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	Token     string
	BotPrefix string

	config *Config
)

type Config struct {
	Token     string `json:"token"`
	BotPrefix string `json:"botPrefix"`
}

// ReadConfig reads the config.json file and unmarshals it into the Config struct
func ReadConfig() error {

	fmt.Println("Reading config.json...")
	file, err := os.ReadFile("./config.json")

	if err != nil {
		return err
	}

	fmt.Println("Unmarshalling config.json...")

	// unmarshall file into config struct
	err = json.Unmarshal(file, &config)

	if err != nil {
		fmt.Println("Error unmarshalling config.json")
		return err
	}

	Token = config.Token
	BotPrefix = config.BotPrefix

	Token = config.Token
	BotPrefix = config.BotPrefix

	return nil
}
```

What happened here:

- We created a struct to hold the config.json data
- We created a function to read the config.json file
- We unmarshalled the config.json file into the Config struct
- We assigned the values to the variables Token and BotPrefix

## 6. bot/bot.go
```go
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

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
	}

	BotId = u.ID

	goBot.AddHandler(messageHandler)

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
```

What happened here:

- We created a function to start the bot
- We created a function to handle messages 

*MessageHandlers* will essentially receive a message and a session, which will allow us to send messages back to the channel depending on the contents of the message.

## 7. Finally: main.go 

Lets bind everything together and start the bot. 

```go
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

```

What happened here:

- We read the config.json file
- We started the bot
- We created a channel to keep the bot running

## 8. Run the bot

```bash
go build
go run main.go
```

## 9. Test the bot

- Go a private conversation with the bot and type `!ping` and the bot should reply with `pong!`
- Type `@botname ping` and the bot should reply with `pong!`

How cool is that?

## 10. Next steps

- Add more commands to the bot. 
- Add a database to store information.