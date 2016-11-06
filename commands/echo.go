package commands

import (
	"github.com/godwhoa/ezbot/ezbot"
)

/*
  Echos back to user, example irc log
  user: bot!
  bot: user!
*/
type Echo struct {
	ezbot.Command
}

func NewEcho(nick string) *Echo {
	echo := new(Echo)
	echo.Pattern = "^" + nick + "!$"
	return echo
}

func (e *Echo) Execute(user string, msg string, args []string) {
	e.SChan <- user + "!"
}
