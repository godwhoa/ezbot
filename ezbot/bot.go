package ezbot

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"time"

	"gopkg.in/sorcix/irc.v1"
)

// Bot status
const (
	CONNECTED = iota
	CONNECTING
	DISCONNECTED
)

// Error if an connection is already open
var ErrAlreadyConnected = errors.New("Bot already connected")

// Config contains Bot configuration
type Config struct {
	Nick    string
	Addr    string
	Channel string
}

// Bot contains Log channel for getting propagated logs
type Bot struct {
	config    Config
	conn      *irc.Conn
	timeout   time.Duration
	commands  []ICommand
	send      chan string
	status    int
	onconnect func()
	Log       chan string
}

// New creates instance of Bot with initial state
func New() *Bot {
	bot := new(Bot)
	bot.status = DISCONNECTED
	bot.timeout = 300 * time.Second
	bot.send = make(chan string)
	bot.Log = make(chan string)

	return bot
}

// Connect to server, init commands, start send and read loops.
func (b *Bot) Connect(config Config) error {
	if b.status != DISCONNECTED {
		return ErrAlreadyConnected
	}
	b.status = CONNECTING
	b.config = config

	netconn, ircconn, err := b.makeIRCConn()
	if err != nil {
		return err
	}
	b.conn = ircconn

	b.Log <- fmt.Sprintf("[bot] connecting addr: %s chan: %s", b.config.Addr, b.config.Channel)
	b.JoinCmds(false)
	b.initcmds()
	go b.sendloop()
	return b.readloop(netconn, ircconn)
}

// Tries to make tls connection if that fails makes regular tcp connection.
func tlsOrtcp(addr string) (net.Conn, error) {
	var err error
	var conn net.Conn
	conn, err = tls.Dial("tcp", addr, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			return conn, err
		}
	}
	return conn, nil
}

// Creates a tls/tcp connection and then creates an irc connection with it
func (b *Bot) makeIRCConn() (net.Conn, *irc.Conn, error) {
	netconn, err := tlsOrtcp(b.config.Addr)
	if err != nil {
		b.Log <- fmt.Sprintf("[bot] connection err: %s", err.Error())
		return netconn, nil, err
	}
	return netconn, irc.NewConn(netconn), nil
}

// Inits commands
func (b *Bot) initcmds() {
	for _, command := range b.commands {
		command.Once()
	}
}

// Sends messages from send channel
func (b *Bot) sendloop() {
	for b.status != DISCONNECTED {
		b.Send(<-b.send)
	}
}

// Reads IRC messages and fires appropriate reponses
func (b *Bot) readloop(netconn net.Conn, ircconn *irc.Conn) error {
	for b.status != DISCONNECTED {
		netconn.SetDeadline(time.Now().Add(b.timeout))
		message, err := ircconn.Decode()
		if err != nil {
			b.Log <- fmt.Sprintf("[bot] decode err: %s", err.Error())
			return err
		}

		switch message.Command {
		case irc.RPL_WELCOME:
			ircconn.Encode(&irc.Message{Command: irc.JOIN,
				Params: []string{b.config.Channel}})
		case irc.ERR_NICKNAMEINUSE:
			b.JoinCmds(true)
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
	return ircconn.Close()
}

// OnConnect sets callback to fire connecting
func (b *Bot) OnConnect(onconnect func()) {
	b.onconnect = onconnect
}

// Disconnect bot
func (b *Bot) Disconnect() {
	b.status = DISCONNECTED
}

// JoinCmds sends commands to set nick and it retries if nick is taken
func (b *Bot) JoinCmds(taken bool) {
	if taken {
		b.config.Nick += "_"
	}
	b.conn.Encode(&irc.Message{Command: irc.NICK,
		Params: []string{b.config.Nick}})
	b.conn.Encode(&irc.Message{Command: irc.USER,
		Params: []string{b.config.Nick, "0", "*", b.config.Nick}})
	b.Log <- fmt.Sprintf("[bot] setnick: %s", b.config.Nick)
}
