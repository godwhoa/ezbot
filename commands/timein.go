package commands

import (
	"fmt"
	"github.com/godwhoa/ezbot/ezbot"
	"strings"
	"time"
)

const (
	layout = "Jan 2, 2006 at 3:04pm -0700"
	help   = "%s: .timein <IANA timezone>"
	ref    = "Ref:https://en.wikipedia.org/wiki/List_of_tz_database_time_zones"
	note   = "Note:You can use aliases for UK,AU,FIN,IN"
)

func alias(arg string) string {
	ori := arg
	arg = strings.ToLower(arg)
	switch arg {
	case "uk":
		return "GB"
	case "india":
		return "Asia/Calcutta"
	case "in":
		return "Asia/Calcutta"
	case "fin":
		return "EET"
	case "finland":
		return "EET"
	case "au":
		return "Australia/Canberra"
	default:
		return ori
	}
}

type TimeIn struct {
	ezbot.Command
}

func NewTimeIn() *TimeIn {
	timein := new(TimeIn)
	timein.Pattern = "^.timein"
	return timein
}

func (t *TimeIn) Execute(user string, msg string, args []string) {
	if len(args) < 2 {
		t.SChan <- fmt.Sprintf(help, user)
		t.SChan <- ref
		t.SChan <- note
	} else {
		zone := alias(args[1])
		loc, err := time.LoadLocation(zone)
		if err != nil {
			t.SChan <- "Invalid zone."
			fmt.Printf("LoadLocation err: %v", err)
			return
		}
		now := time.Now().In(loc)
		t.SChan <- fmt.Sprintf("%s : %s", zone, now.Format(layout))
	}
}
