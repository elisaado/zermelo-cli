# zermelo-cli
Unofficial command line interface application to access Zermelo (zportal)

## Usage
When no arguments are provided to the executable, it will start an interactive REPL-ish environment where one can execute multiple commands

### Commands
\* Optional

|Command|Parameters|Default action (without parameters)|Description|
|-------|----------|-----------------------------------|-----------|
|help|Command name*|Shows this help|Shows help for a specific command|
|init|Organisation, authentication code|Starts an interactive wizardish thing where the user can enter their organisation and authentication code, then initializes (fetches authentication token) the CLI|Initializes (fetches the authentication token) the CLI|
|show|Day* (parameters like: today, tomorrow or integers where 0 is today, 1 is tomorrow and 6 is today but next week)|Shows schedule for today|Shows schedule for specific day|
|me|None||Shows all info Zermelo knows about you|
|info|None||Shows info about zermelo-cli (version, creator)|

## Author
Eli Saado and contributors
