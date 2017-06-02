package commands

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/godwhoa/ezbot/ezbot"
)

/*
	user: http://exez.in
	ezbot: exezin
*/
type Title struct {
	ezbot.Command
}

func NewTitle() *Title {
	title := &Title{}
	title.Pattern = `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`
	return title
}

func (t *Title) OnMsg(user string, msg string) {
	urls := t.Reg.FindAllString(msg, -1)
	for _, url := range urls {
		if url != "" {
			doc, err := goquery.NewDocument(url)
			if err != nil {
				return
			}
			t.Log <- fmt.Sprintf("[title] url: %s", url)
			t.SChan <- doc.Find("title").First().Text()
		}
	}
}
