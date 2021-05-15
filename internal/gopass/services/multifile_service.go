package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/viper"
	"gopass/gopass/internal/gopass/interfaces"
)

// Manager manages working with multiple files
type Manager struct {
	Keepass   interfaces.KeepassService
	Prompt    interfaces.PromptService
	Multifile interfaces.MultiFileService
}

type multifile struct {
	config  config
	keepass interfaces.KeepassService
}

type config struct {
	Locations []string
}

var configName = "gopass"

func (data *multifile) nextToOpen() string {
	for _, location := range data.config.Locations {
		if !data.keepass.IsOpen(location) {
			return location
		}
	}
	return ""
}

// Open opens all locations from the config
func (data *multifile) Open() {
	viper.SetConfigName(configName)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	data.config.Locations = viper.GetStringSlice("locations")

	var location string = data.nextToOpen()
	var index = data.keepass.DBCount()
	for location != "" {
		data.keepass.Open(location)
		index++
		if index != data.keepass.DBCount() {
			break
		}
		location = data.nextToOpen()
	}
}

// ContinueOpen continues opening remaining kdbx files from the config, which have not been opened yet.
// If the parameter file is passed in, method verifies if the file has been already opened or not and opens it if needed.
// This behavior is desired when manually opening the databases and opening them from the config.
func (data *multifile) ContinueOpen(file string) {
	if file != "" && !data.keepass.IsOpen(file) {
		data.keepass.Open(file)
		return
	}
	data.Open()
}

// InitConfig initializes the configuration file in user's home folder.
func (data *multifile) InitConfig() {
	var directory = os.ExpandEnv("$HOME/.gopass")
	var destinationFile = os.ExpandEnv("$HOME/.gopass/" + configName + ".json")

	if _, err := os.Stat(destinationFile); !os.IsNotExist(err) {
		fmt.Printf("Config file already exists: %s\n", destinationFile)
		fmt.Println("Aborting new config creation.")
		return
	}

	if err := os.MkdirAll(directory, 0755); err != nil {
		log.Fatal(err)
	}

	blankConfig := config{
		Locations: []string{
			"/set/path/to/your/keepass.kdbx",
		},
	}

	input, err := json.MarshalIndent(blankConfig, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(destinationFile, input, 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Config initialized in %s\n", destinationFile)
	fmt.Println("Feel free to edit it and add paths to your kdbx database files!")

}

func (data *multifile) LocationsCount() int {
	return len(data.config.Locations)
}

// NewMultifileService returns a new instance of the service
func NewMultifileService(keepass interfaces.KeepassService) interfaces.MultiFileService {
	var multifileInstance = new(multifile)
	multifileInstance.keepass = keepass
	var multifileService interfaces.MultiFileService = multifileInstance
	return multifileService
}

// NewMultifileManager creates and returns the manager
func NewMultifileManager(keepass interfaces.KeepassService, prompt interfaces.PromptService, multifile interfaces.MultiFileService) *Manager {
	var managerInstance = new(Manager)
	managerInstance.Keepass = keepass
	managerInstance.Prompt = prompt
	managerInstance.Multifile = multifile
	return managerInstance
}
