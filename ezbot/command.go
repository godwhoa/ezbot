package ezbot

import (
	"regexp"
)

type Command struct {
	Pattern string
	Reg     *regexp.Regexp
	SChan   chan string
}

type ICommand interface {
	Init(schan chan string)
	Match(msg string) bool
	Once()
	Execute(user string, msg string, args []string)
	OnJoin(user string)
	OnLeave(user string)
	OnMsg(user string, msg string)
}

// Sets channel with one passed by Bot
// Compiles regex for performance
func (c *Command) Init(schan chan string) {
	c.SChan = schan
	c.Reg = regexp.MustCompile(c.Pattern)
}

// Getting around interface restriction
func (c *Command) Match(msg string) bool {
	return c.Reg.MatchString(msg)
}

// Only executes once
func (c *Command) Once() {
}

// Executes when it matches with the regex.
func (c *Command) Execute(user string, msg string, args []string) {
}

// Executes when a user joins channel
func (c *Command) OnJoin(user string) {
}

// Executes when a user leaves channel
func (c *Command) OnLeave(user string) {
}

// Executes when user sends a message
func (c *Command) OnMsg(user string, msg string) {
}
