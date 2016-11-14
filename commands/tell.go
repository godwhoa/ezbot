package commands

import (
	"fmt"
	"github.com/godwhoa/ezbot/ezbot"
	"github.com/hako/durafmt"
	"strings"
	"time"
)

type Letter struct {
	from string
	body string
	when time.Time
}

/*
  Acts as mailbox for offline users
  from: .tell to <letter>
  bot: message
  bot: 1hr ago
*/
type Tell struct {
	ezbot.Command
	mailbox map[string][]Letter
}

func NewTell() *Tell {
	tell := &Tell{}
	tell.Pattern = "^.tell"
	tell.mailbox = make(map[string][]Letter)
	return tell
}

// Notifies when user has some letters
func (t *Tell) Notify(user string) {
	if letters, ok := t.mailbox[user]; ok {
		for _, letter := range letters {
			time_since := time.Since(letter.when)
			t.SChan <- fmt.Sprintf("%s, %s left this message for you: %s", user, letter.from, letter.body)
			t.SChan <- fmt.Sprintf("%s ago", durafmt.Parse(time_since).String())
		}
		delete(t.mailbox, user)
	}
}

func (t *Tell) Execute(user string, msg string, args []string) {
	if len(args) < 2 {
		t.SChan <- "Usage: .tell <to> <msg>"
		return
	}
	from := user
	to := args[1]
	body := strings.Join(args[2:], " ")
	t.mailbox[to] = append(t.mailbox[to], Letter{from, body, time.Now()})
	t.SChan <- "Okie Dokie!"
}

func (t *Tell) OnJoin(user string) {
	t.Notify(user)
}

func (t *Tell) OnMsg(user string, msg string) {
	t.Notify(user)
}
