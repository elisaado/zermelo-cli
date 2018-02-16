package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 || (os.Args[1] == "help" && len(os.Args) < 3) {
		fmt.Println(GetHelp())
	}
}

func GetHelp() string {
	return `zermelo [command] [command args]
  Commands:
    help                            - Show this help
    help [Command name]             - Show help for a specific command
    init                            - Show an interactive prompt asking for your organisation and authentication code
    init [organisation] [auth code] - Initialize Zermelo CLI
    show                            - Show schedule for today
    show [day]                      - Show schedule for specific day
    me                              - Show all info Zermelo knows about you
    info                            - Show info about zermelo-cli (version, author)
`
}
