ezbot
==========
An IRC framework in Go aimed for ease of use and flexibility.


## Usage

### CLI from linux binary
```
./ezbot_linux
```

### CLI from source
```
go get -v github.com/godwhoa/ezbot
cd $GOPATH/src/github.com/godwhoa/ezbot
go run main.go
```

### Package
```go
package main
import(
	"github.com/godwhoa/ezbot/ezbot"
)

type Echo struct {
	ezbot.Command
}

func NewEcho(nick string) *Echo {
	echo := new(Echo)
	echo.Pattern = "^" + nick + "!$" // regex for the command
	return echo
}

// Runs when it matches with the regex
func (e *Echo) Execute(user string, msg string, args []string) {
	e.SChan <- user + "!"
}

// Other methods like OnJoin, OnLeave, OnMsg are also available
// Note: you can leave them out if you wish to
// Also see commands dir. for examples and go read thru ezbot/command.go

func main(){
	// create a bot
	m_bot := ezbot.New("ezbot","#ezbot","chat.freenode.net:6697")
	// Init command
	echo := NewEcho("ezbot")
	// Add commands
	m_bot.Add(echo)
	// Connects and spins-off a read/write loop
	m_bot.Connect()
}
```


## TODO
 + ~~Add more commands~~
 + ~~Fix url regex~~
 + ~~Handle TLS servers~~
 + ~~Config file~~
 + Handle nick collision
 + Better naming of types, variables etc.
 + Add tests