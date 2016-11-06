package commands

import (
	"encoding/json"
	"fmt"
	"github.com/godwhoa/ezbot/ezbot"
	"net"
	"net/http"
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
			fmt.Println("Failed to parse json.")
		}
		g.SChan <- fmt.Sprintf("New commit: %s", parsed.Commits[0].URL)
		fmt.Fprintf(w, "OK.")
	}
}

func (g *Git) Once() {
	go func() {
		fmt.Printf("webhook port: %s\n", g.port)
		http.HandleFunc("/endpoint", g.Endpoint)
		http.ListenAndServe(g.port, nil)
	}()
}
