package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/godwhoa/ezbot/commands"
	"github.com/godwhoa/ezbot/ezbot"
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
	fmt.Printf("Commands: ")
	for _, cmd := range config.Commands {
		fmt.Printf("%s ", cmd)
	}
	fmt.Printf("\n")

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
		case "reminder":
			m_bot.AddCmd(commands.NewReminder())
		}
	}
	go func() {
		for {
			fmt.Println(<-m_bot.Log)
		}
	}()
	time.AfterFunc(time.Minute*2, func() {
		m_bot.Disconnect()
	})
	fmt.Println(m_bot.Connect())
}
