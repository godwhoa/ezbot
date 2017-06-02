package commands

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/godwhoa/ezbot/ezbot"
)

func GetPort() string {
	l, _ := net.Listen("tcp", ":0")
	defer l.Close()
	return l.Addr().String()
}

type Payload struct {
	Commits []struct {
		URL string `json:"url"`
	}
}

type Git struct {
	ezbot.Command
	port string
}

func NewGit() *Git {
	git := new(Git)
	git.port = GetPort()
	return git
}

func (g *Git) Endpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var parsed Payload
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&parsed)
		if err != nil {
			g.Log <- fmt.Sprintf("[git] err: %s", err.Error())
		} else {
			g.Log <- fmt.Sprintf("[git] notified commit url: %s", parsed.Commits[0].URL)
			g.SChan <- fmt.Sprintf("New commit: %s", parsed.Commits[0].URL)
			w.Write([]byte("OK."))
		}
	}
}

func (g *Git) Once() {
	go func() {
		fmt.Printf("webhook port: %s\n", g.port)
		http.HandleFunc("/endpoint", g.Endpoint)
		http.ListenAndServe(g.port, nil)
	}()
}
