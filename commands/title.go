package commands

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/godwhoa/ezbot/ezbot"
)

/*

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
	url := t.Reg.FindString(msg)
	if url != "" {
		doc, err := goquery.NewDocument(url)
		if err != nil {
			return
		}
		t.SChan <- doc.Find("title").Text()

	}
}
