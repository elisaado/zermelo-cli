package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 || (os.Args[1] == "help" && len(os.Args) < 3) {
		fmt.Println(getHelp())
		os.Exit(1)
	}

	switch os.Args[1] {
	case "help":
		fmt.Println(getHelpFor(os.Args[2]))
	default:
		fmt.Println(getHelp())
	}
}

func getHelp() string {
	return `zermelo-cli is an unofficial command line interface application to access Zermelo (zportal)

Usage:
  zermelo [command] [command args]
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

func getHelpFor(command string) string {
	switch command {
	case "init":
		return `Init is used to initialize zermelo-cli (get the authentication token), when no arguments are provided, zermelo-cli will start an interactive prompt where it will ask the user for their organisation and authentication code, and then it will fetch the authentication token used for further requests. It is saved in plain text in a json file to ~/.config/zermelo-cli/config.json, otherwise it will do the same thing without the interactive prompt`
	case "show":
		return "Show is used to show the schedule for a day (kind of the core of this program), when no arguments are provided it will show the schedule for today. Possible arguments are: today, tomorrow, and integers (where 0 is today, 1 is tomorrow, 6 is next week, etc.)"
	case "me":
		return "Me is used to see who's currently logged in, it returns all info Zermelo knows about you (first name, last name, etc.)"
	case "info":
		return "Info is used to get version (and author :D) info, useful for debugging"
	default:
		return getHelp()
	}
}
