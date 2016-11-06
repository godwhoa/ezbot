package commands

import (
	"fmt"
	"github.com/godwhoa/ezbot/ezbot"
	"github.com/hako/durafmt"
	"time"
)

/*
  Sends last time a user was active
  user: .seen user2
  bot: user2 was last seen 2 hours ago
*/
type Seen struct {
	ezbot.Command
	store map[string]time.Time
}

func NewSeen() *Seen {
	seen := &Seen{}
	seen.Pattern = "^.seen"
	seen.store = make(map[string]time.Time)
	return seen
}

func (s *Seen) Execute(user string, msg string, args []string) {
	user2 := args[1]
	if last_seen, ok := s.store[user2]; ok {
		time_since := time.Since(last_seen)
		pretty_time := durafmt.Parse(time_since).String()
		s.SChan <- fmt.Sprintf("%s was last seen %s ago.", user2, pretty_time)
	} else {
		s.SChan <- fmt.Sprintf("No log for %s found.", user2)
	}
}

func (s *Seen) OnJoin(user string) {
	s.store[user] = time.Now()
}

func (s *Seen) OnLeave(user string) {
	s.store[user] = time.Now()
}

func (s *Seen) OnMsg(user string, msg string) {
	s.store[user] = time.Now()
}
