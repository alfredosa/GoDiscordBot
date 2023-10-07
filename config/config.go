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
