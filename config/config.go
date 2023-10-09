package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	Token     string
	BotPrefix string
	ChannelID string
	GuildId   string

	config *Config
)

type Config struct {
	Token     string `json:"token"`
	BotPrefix string `json:"botPrefix"`
	ChannelID string `json:"channelID"`
	GuildId   string `json:"guildId"`
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
	ChannelID = config.ChannelID
	GuildId = config.GuildId

	return nil
}
