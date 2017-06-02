package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/godwhoa/ezbot/commands"
	"github.com/godwhoa/ezbot/ezbot"
)

type configFile struct {
	Nick     string   `json:"nick"`
	Channel  string   `json:"channel"`
	Addr     string   `json:"addr"`
	Commands []string `json:"commands"`
}

func (c *configFile) print() {
	fmt.Printf("Nick: %s Addr: %s Channel: %s", c.Nick, c.Addr, c.Channel)
	fmt.Printf("Commands: ")
	for _, command := range c.Commands {
		fmt.Printf("%s ", command)
	}
	fmt.Printf("\n")
}

func (c *configFile) toBotConfig() ezbot.Config {
	return ezbot.Config{Nick: c.Nick, Channel: c.Channel, Addr: c.Addr}
}

func loadConfigFile(path string) (configFile, error) {
	var config configFile

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(file, &config); err != nil {
		return config, err
	}

	return config, nil
}

func main() {
	configfile, err := loadConfigFile("./config.json")
	if err != nil {
		log.Fatalf("Failed to load configfile")
	}

	bot := ezbot.New()
	for _, c := range configfile.Commands {
		switch c {
		case "echo":
			bot.AddCmd(commands.NewEcho(configfile.Nick))
		case "seen":
			bot.AddCmd(commands.NewSeen())
		case "tell":
			bot.AddCmd(commands.NewTell())
		case "timein":
			bot.AddCmd(commands.NewTimeIn())
		case "title":
			bot.AddCmd(commands.NewTitle())
		case "git":
			bot.AddCmd(commands.NewGit())
		case "reminder":
			bot.AddCmd(commands.NewReminder())
		}
	}
	go func() {
		for {
			fmt.Println(<-bot.Log)
		}
	}()
	log.Fatal(bot.Connect(configfile.toBotConfig()))
}
