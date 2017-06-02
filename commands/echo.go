package commands

import (
	"fmt"

	"github.com/godwhoa/ezbot/ezbot"
)

/*
  Echos back to user, example irc log
  user: bot!
  bot: user!
*/
type Echo struct {
	ezbot.Command
	nick string
}

func NewEcho(nick string) *Echo {
	echo := new(Echo)
	echo.Pattern = "^" + nick + "!$"
	echo.nick = nick
	return echo
}

func (e *Echo) Execute(user string, msg string, args []string) {

	if user == "exezin" {
		e.SChan <- e.nick + " <3 " + user
	} else {
		e.SChan <- user + "!"
	}
	e.Log <- fmt.Sprintf("[echo] echoed back to %s", user)
}
