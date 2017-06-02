package ezbot

import (
	"fmt"
	"strings"

	"gopkg.in/sorcix/irc.v1"
)

// Adds commands to be executed
func (b *Bot) AddCmd(commands ...ICommand) {
	for _, command := range commands {
		command.Init(b.send, b.Log)
		b.commands = append(b.commands, command)
	}
}

// Send text message
func (b *Bot) Send(msg string) {
	b.conn.Encode(&irc.Message{
		Command:  irc.PRIVMSG,
		Params:   []string{b.config.Channel},
		Trailing: msg,
	})
}

// Channel joins
func (b *Bot) Join(message *irc.Message) {
	nick := message.Name
	if nick != b.config.Nick {
		for _, command := range b.commands {
			command.OnJoin(nick)
		}
	} else {
		b.status = CONNECTED
		if b.onconnect != nil {
			b.onconnect()
		}
		b.Log <- fmt.Sprintf("[bot] joined: %s/%s", b.config.Addr, b.config.Channel)
	}
}

// Channel leaves
func (b *Bot) Leave(message *irc.Message) {
	nick := message.Name
	for _, command := range b.commands {
		command.OnLeave(nick)
	}
}

// Channel messages
func (b *Bot) Msg(message *irc.Message) {
	nick := message.Name
	msg := message.Trailing
	arg := strings.Split(msg, " ")
	for _, command := range b.commands {
		if command.Match(msg) {
			command.Execute(nick, msg, arg)
		}
		command.OnMsg(nick, msg)
	}
}

// PONG reply
func (b *Bot) Pong(message *irc.Message) {
	b.conn.Encode(&irc.Message{
		Command:  irc.PONG,
		Params:   message.Params,
		Trailing: message.Trailing,
	})
}
