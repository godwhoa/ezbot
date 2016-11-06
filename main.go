package main

import (
	"encoding/json"
	"fmt"
	"github.com/godwhoa/ezbot/commands"
	"github.com/godwhoa/ezbot/ezbot"
	"io/ioutil"
)

var (
	m_nick    = "ezbot"
	m_channel = "#ezirc"
	m_addr    = "chat.freenode.net:6667"
)

type Config struct {
	Nick     string   `json:"nick"`
	Channel  string   `json:"channel"`
	Addr     string   `json:"addr"`
	Commands []string `json:"commands"`
}

func main() {

	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
	}
	var config Config
	json.Unmarshal(file, &config)
	m_bot := ezbot.New(config.Nick, config.Channel, config.Addr)
	for _, c := range config.Commands {
		switch c {
		case "echo":
			m_bot.AddCmd(commands.NewEcho(config.Nick))
		case "seen":
			m_bot.AddCmd(commands.NewSeen())
		case "tell":
			m_bot.AddCmd(commands.NewTell())
		case "timein":
			m_bot.AddCmd(commands.NewTimeIn())
		case "title":
			m_bot.AddCmd(commands.NewTitle())
		case "git":
			m_bot.AddCmd(commands.NewGit())
		}
	}

	m_bot.Connect()
}
