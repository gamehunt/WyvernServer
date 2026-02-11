package client

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
)

var serverUrl string

func execExitCommand(_ []string) error {
	fmt.Println("Bye")
	os.Exit(0)
	return nil
}

func execHealthCommand(_ []string) error {
	err := checkHealth(serverUrl)	
	if err != nil {
		fmt.Println("Health check failed")
		return err
	} else {
		fmt.Println("Health check success")
		return nil
	}
}

type CommandFunc func(args []string) error

type Command struct {
	Description string
	Exec        CommandFunc
}

var commands = map[string]Command {
	"health": {
		Description: "Check Health",
		Exec: execHealthCommand,
	},
	"exit": {
		Description: "Exit",
		Exec: execExitCommand,
	},
}

func makeSuggestions() []prompt.Suggest {
	suggestions := []prompt.Suggest{}
	for name, cmd := range commands {
		suggestions = append(suggestions, prompt.Suggest{Text: name, Description: cmd.Description})	
	}
	return suggestions
}

var suggestions = makeSuggestions()

func executor(cmd string) {
	words := strings.Fields(cmd)

	command := words[0]
	args    := words[1:]

	rec, ok := commands[command]

	if !ok {
		fmt.Println("No such command")	
		return
	}

	err := rec.Exec(args)

	if err != nil {
		fmt.Println("An error occured during command execution: %v", err)
	}
}

func completer(in prompt.Document) []prompt.Suggest {
	w := in.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}
	return prompt.FilterHasPrefix(suggestions, w, true)
}

func main() {

	serverUrlPtr := flag.String("server", "http://localhost:5000", "target server url")
	flag.Parse()

	serverUrl = *serverUrlPtr

	log.Printf("Connecting to %s", serverUrl)

	err := checkHealth(serverUrl)

	if err != nil {
		log.Fatalf("Initial check health failed: %v", err)
	}

	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("> "),
	)
	p.Run()
}
