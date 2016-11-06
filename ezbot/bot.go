package ezbot

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/sorcix/irc.v1"
	"net"
	"strings"
)

type Bot struct {
	conn       *irc.Conn
	tls_config *tls.Config
	addr       string
	nick       string
	channel    string
	commands   []ICommand
	schan      chan string
}

func New(nick string, channel string, addr string) *Bot {
	bot := new(Bot)
	bot.nick = nick
	bot.channel = channel
	bot.addr = addr
	bot.schan = make(chan string)
	bot.tls_config = &tls.Config{InsecureSkipVerify: true}
	return bot
}

// Make a conn. and spin-off read/write loops.
func (b *Bot) Connect() {
	tconn, err := tls.Dial("tcp", b.addr, b.tls_config)
	if err != nil {
		conn, err := net.Dial("tcp", b.addr)
		if err != nil {
			fmt.Println("Failed to connect.")
		}
		b.conn = irc.NewConn(conn)
	} else {
		b.conn = irc.NewConn(tconn)
	}

	fmt.Printf("Connected to %s\nJoining %s\n", b.addr, b.channel)

	b.JoinCmds()
	for _, command := range b.commands {
		command.Once()
	}
	// Send loop
	go func() {
		for {
			b.Send(<-b.schan)
		}
	}()

	// Read loop
	for {
		message, err := b.conn.Decode()
		if err != nil {
			fmt.Println("Failed to decode.")
			break
		}

		switch message.Command {
		case irc.JOIN:
			b.Join(message)
		case irc.PART:
			b.Leave(message)
		case irc.PRIVMSG:
			b.Msg(message)
		case irc.PING:
			b.Pong(message)
		}
	}

	b.conn.Close()
}

// Sets nick and joins channel
func (b *Bot) JoinCmds() {
	b.conn.Encode(&irc.Message{Command: irc.NICK,
		Params: []string{b.nick}})
	b.conn.Encode(&irc.Message{Command: irc.USER,
		Params: []string{b.nick, "0", "*", b.nick}})
	b.conn.Encode(&irc.Message{Command: irc.JOIN,
		Params: []string{b.channel}})
}

// Adds commands to be executed
func (b *Bot) AddCmd(commands ...ICommand) {
	for _, command := range commands {
		command.Init(b.schan)
		b.commands = append(b.commands, command)
	}
}

/* Response methods */
// Send text message
func (b *Bot) Send(msg string) {
	b.conn.Encode(&irc.Message{
		Command: irc.PRIVMSG,
		Params:  []string{b.channel, ":" + msg},
	})
}

// Channel joins
func (b *Bot) Join(message *irc.Message) {
	nick := message.Name
	if nick != b.nick {
		for _, command := range b.commands {
			command.OnJoin(nick)
		}
	} else {
		fmt.Printf("Joined %s\n", b.channel)
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

// PING reply
func (b *Bot) Pong(message *irc.Message) {
	b.conn.Encode(&irc.Message{
		Command:  irc.PONG,
		Params:   message.Params,
		Trailing: message.Trailing,
	})
}
