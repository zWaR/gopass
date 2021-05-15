package services

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	"gopass/gopass/internal/gopass/interfaces"
)

// Prompt is the basic data type and receiver.
type gopassPrompt struct {
	keepass   interfaces.KeepassService
	multifile interfaces.MultiFileService
	dbFile    string
}

var suggestions = []prompt.Suggest{
	{Text: "u", Description: "Copy username"},
	{Text: "p", Description: "Copy password"},
	{Text: "find", Description: "Find the entry in the keepass file"},
	{Text: "cont", Description: "Continue to open files from config"},
	{Text: "ls", Description: "Lists currently opened KeePass databases"},
	{Text: "open", Description: "Open a KeePass database"},
	{Text: "close", Description: "Close a KeePass Database"},
	{Text: "help", Description: "Show usage help"},
	{Text: "quit", Description: "Quit gopass"},
}

func completer(in prompt.Document) []prompt.Suggest {
	word := in.GetWordBeforeCursor()
	if word == "" {
		return []prompt.Suggest{}
	}
	return prompt.FilterHasPrefix(suggestions, word, true)
}

func getParameter(blocks []string) string {
	if len(blocks) == 1 {
		return ""
	}

	var parameter string
	index := 1
	for index < len(blocks) {
		parameter += " " + blocks[index]
		index++
	}
	return strings.TrimPrefix(parameter, " ")

}

func showHelp() {
	for _, suggestion := range suggestions {
		fmt.Printf("%s\t%s\n", suggestion.Text, suggestion.Description)
	}
	fmt.Println()
}

func (p *gopassPrompt) executor(in string) {
	in = strings.TrimSpace(in)
	var digitCheck = regexp.MustCompile(`^[0-9]+$`)

	blocks := strings.Split(in, " ")
	parameter := getParameter(blocks)
	switch blocks[0] {
	case "find":
		p.keepass.Find(parameter)
	case "u":
		if digitCheck.MatchString(parameter) {
			id, _ := strconv.Atoi(parameter)
			p.keepass.Copy(id, "username")
		} else if p.keepass.Find(parameter) == 1 {
			p.keepass.Copy(0, "username")
		}
	case "p":
		if digitCheck.MatchString(parameter) {
			id, _ := strconv.Atoi(parameter)
			p.keepass.Copy(id, "password")
		} else if p.keepass.Find(parameter) == 1 {
			p.keepass.Copy(0, "password")
		}
	case "open":
		if parameter != "" {
			p.dbFile = parameter
			p.keepass.Open(parameter)
		}
	case "close":
		if digitCheck.MatchString(parameter) {
			id, _ := strconv.Atoi(parameter)
			p.keepass.Close(id)
		}
	case "cont":
		p.multifile.ContinueOpen(p.dbFile)
		p.dbFile = ""
	case "ls":
		p.keepass.PrintDatabases()
	case "help":
		showHelp()
	case "quit":
		os.Exit(0)
	default:
		if blocks[0] != "" {
			fmt.Println("Unrecognized command. Run 'help' for more information.")
		}
	}
}

// Start starts a new command prompt.
func (p *gopassPrompt) Start() {
	var prompt = prompt.New(
		p.executor,
		completer,
		prompt.OptionPrefix("> "),
		prompt.OptionTitle("gopass-prompt"),
	)
	prompt.Run()
}

// NewPromptService is a PromptService provider.
func NewPromptService(keepass interfaces.KeepassService, multifile interfaces.MultiFileService) interfaces.PromptService {
	var gopassPromptInstance = new(gopassPrompt)
	gopassPromptInstance.keepass = keepass
	gopassPromptInstance.multifile = multifile
	gopassPromptInstance.dbFile = ""
	var prompt interfaces.PromptService = gopassPromptInstance
	return prompt
}
