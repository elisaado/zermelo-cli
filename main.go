// This file is the file the user is going to interact with
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/shibukawa/configdir"
)

type Config struct {
	Organisation string `json:"organisation"`
	Token        string `json:"token"`
}

var config Config

var helpString string = `zermelo-cli is an unofficial command line interface application to access Zermelo (zportal)

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

func main() {
	// Read config file
	configDirs := configdir.New("elisaado", "zermelo-cli")
	folder := configDirs.QueryFolderContainsFile("config.json")
	if folder != nil {
		data, _ := folder.ReadFile("config.json")
		json.Unmarshal(data, &config)
	} else {
		fmt.Println("No config found... Initializing zermelo-cli...")
		initialize()
		os.Exit(0)
	}

	// check if there are no args or if the second command is help and it is not help [command]
	if len(os.Args) < 2 || (os.Args[1] == "help" && len(os.Args) < 3) {
		fmt.Println(helpString)
		os.Exit(1)
	}

	baseurl = "https://" + config.Organisation + ".zportal.nl/api/v3"

	// Check which command it actually is
	switch os.Args[1] {
	case "help":
		fmt.Println(getHelpFor(os.Args[2]))
	case "init":
		fmt.Println("No need to reinitialize... Delete ~/.config/elisaado/zermelo-cli/config.json to log out")
	case "show":
		var appointments []Appointment
		if len(os.Args) < 3 || os.Args[2] == "today" || os.Args[2] == "0" {
			appointments = fetchAppointments(config.Organisation, config.Token, int(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).Unix()), int(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).Unix()+60*60*24))
		} else if os.Args[2] == "tomorrow" {
			appointments = fetchAppointments(config.Organisation, config.Token, int(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 0, 0, 0, 0, time.Local).Unix()), int(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 0, 0, 0, 0, time.Local).Unix()+60*60*24))
		} else {
			day, err := strconv.Atoi(os.Args[2])
			if err != nil || os.Args[2] != strconv.Itoa(day) {
				fmt.Println("Invalid day")
				os.Exit(1)
			}
			appointments = fetchAppointments(config.Organisation, config.Token, int(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+day, 0, 0, 0, 0, time.Local).Unix()), int(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+day, 0, 0, 0, 0, time.Local).Unix()+60*60*24))
		}

		// Sort the appointments
		sort.Slice(appointments, func(i, j int) bool { return appointments[i].Start < appointments[j].Start })

		// Print the appointments in a nice table
		fmt.Println(appointmentPrint(appointments))

	default:
		fmt.Println(helpString)
	}
}

func getHelpFor(command string) string {
	// Check which command it is the user needs help with
	switch command {
	case "init":
		return `Init is used to initialize zermelo-cli (get the authentication token), when no arguments are provided, zermelo-cli will start an interactive prompt where it will ask the user for their organisation and authentication code, and then it will fetch the authentication token used for further requests. It is saved in plain text in json format to ~/.config/elisaado/zermelo-cli/config.json`
	case "show":
		return "Show is used to show the schedule for a day (kind of the core of this program), when no arguments are provided it will show the schedule for today. Possible arguments are: today, tomorrow, and integers (where 0 is today, 1 is tomorrow, 6 is next week, etc.)"
	case "me":
		return "Me is used to see who's currently logged in, it returns all info Zermelo knows about you (first name, last name, etc.)"
	case "info":
		return "Info is used to get version (and author :D) info, useful for debugging"
	default:
		return helpString
	}
}

func initialize() {
	// Create scanner to read from stdin
	scanner := bufio.NewScanner(os.Stdin)

	// Read organisation from stdin
	fmt.Print("\nPlease enter your organisation: ")
	scanner.Scan()
	organisation := scanner.Text()

	// Read auth code from stdin
	fmt.Print("Please enter your authentication code: ")
	scanner.Scan()
	codeS := strings.Replace(scanner.Text(), " ", "", -1)

	// Show them to the user again
	fmt.Printf("\nOrganisation: %s, authentication code: %s\n\n", organisation, codeS)

	// Convert code to int and check if it's valid
	code, err := strconv.Atoi(codeS)
	if len(codeS) != 12 || err != nil || code == 0 {
		fmt.Println("Invalid authentication code")
		os.Exit(1)
	}

	// Fetch auth token
	fmt.Println("Fetching authentication token...")
	token := fetchAuthToken(organisation, code)
	if token == "" {
		fmt.Println("An error occured... Have you typed everything correctly?")
	}

	fmt.Println("Finished fetching token... Writing it to ~/.config/elisaado/zermelo-cli/config.json...")

	// Write token and organisation to config file for further requests
	configDirs := configdir.New("elisaado", "zermelo-cli")

	config.Organisation = organisation
	config.Token = token

	data, _ := json.Marshal(&config)
	folders := configDirs.QueryFolders(configdir.Global)
	folders[0].WriteFile("config.json", data)
}

func appointmentPrint(appointments []Appointment) string {
	if len(appointments) == 0 {
		return "Nothing."
	}

	// Initialize table
	table := simpletable.New()
	cells := []*simpletable.Cell{}

	// Fill table header with time of lessons
	for _, appointment := range appointments {
		cells = append(cells, &simpletable.Cell{Align: simpletable.AlignCenter, Text: time.Unix(int64(appointment.Start), 0).Format("15:04") + "-" + time.Unix(int64(appointment.End), 0).Format("15:04")})
	}
	table.Header = &simpletable.Header{
		Cells: cells,
	}

	// Fill the body wiith subjects and teachers
	var subjects []*simpletable.Cell
	var teachers []*simpletable.Cell
	for _, appointment := range appointments {
		subjects = append(subjects, &simpletable.Cell{
			Align: simpletable.AlignRight, Text: strings.Join(appointment.Subjects, " ,"),
		})
		teachers = append(teachers, &simpletable.Cell{
			Align: simpletable.AlignRight, Text: strings.Join(appointment.Teachers, ", "),
		})
	}
	table.Body.Cells = append(table.Body.Cells, subjects)
	table.Body.Cells = append(table.Body.Cells, teachers)

	// "Render" table
	table.SetStyle(simpletable.StyleUnicode)
	return table.String()
}
