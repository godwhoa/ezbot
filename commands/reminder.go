package commands

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/godwhoa/ezbot/ezbot"
)

/* Time format parsing */
var ParseErr = errors.New("Illegal syntax")

func Parse(format string) (map[string]int, error) {
	umap := make(map[string]int)

	nums := []int{}
	units := []string{}
	reg := regexp.MustCompile(`([0-9]+)|day|hour|min|sec|d|h|m|s`)
	matches := reg.FindAllString(format, -1)
	for _, match := range matches {
		if i, err := strconv.Atoi(match); err == nil {
			nums = append(nums, i)
		} else if match == "day" || match == "d" {
			units = append(units, "day")
		} else if match == "hour" || match == "h" {
			units = append(units, "hour")
		} else if match == "min" || match == "m" {
			units = append(units, "min")
		} else if match == "sec" || match == "s" {
			units = append(units, "sec")
		} else {
			return umap, ParseErr
		}
	}

	if len(nums) != len(units) {
		return umap, ParseErr
	}
	for i := 0; i < len(nums); i++ {
		umap[units[i]] = nums[i]
	}
	return umap, nil
}

func toDur(umap map[string]int) time.Duration {
	dur := 0 * time.Second
	for u, d := range umap {
		switch u {
		case "day":
			dur += time.Hour * time.Duration(24*d)
		case "hour":
			dur += time.Hour * time.Duration(d)
		case "min":
			dur += time.Minute * time.Duration(d)
		case "sec":
			dur += time.Second * time.Duration(d)
		}
	}
	return dur
}

/*
  Sends you a reminder
  user: .reminder 30min2sec blah blah
  bot: Okie Dokie!
*/
type Reminder struct {
	ezbot.Command
}

func NewReminder() *Reminder {
	reminder := new(Reminder)
	reminder.Pattern = "^.reminder "
	return reminder
}

func (r *Reminder) Execute(user string, msg string, args []string) {
	if len(args) < 2 {
		r.SChan <- "Usage: .reminder <time> <msg>"
		r.SChan <- "time format: 1day2hour3min1sec or 1d2h3m1s or 1d1min"
		r.SChan <- "No restrictions on: mixing up long/short form, providing all units and order of units"
		return
	}
	timeformat := args[1]
	body := strings.Join(args[2:], " ")
	go func() {
		unitmap, err := Parse(timeformat)
		if err != nil {
			r.SChan <- "Incorrect time format"
			r.SChan <- "time format: 1day2hour3min1sec or 1d2h3m1s or 1d1min"
			r.SChan <- "No restrictions on: mixing up long/short form, providing all units and order of units"
			return
		}
		duration := toDur(unitmap)
		r.Log <- fmt.Sprintf("[reminder] %s added reminder for %s", user, timeformat)
		time.AfterFunc(duration, func() {
			r.SChan <- fmt.Sprintf("Reminder from %s ago: %s", timeformat, body)
		})

	}()
}
