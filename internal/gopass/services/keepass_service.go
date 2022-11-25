package services

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"gopass/gopass/internal/gopass/interfaces"

	"github.com/atotto/clipboard"
	"github.com/tobischo/gokeepasslib/v3"
	"golang.org/x/crypto/ssh/terminal"
)

type keepass struct {
	Filepaths   []string
	DBs         []*gokeepasslib.Database
	lastResults []result
	DbNames     []string
}

type result struct {
	entry       gokeepasslib.Entry
	parentGroup string
}

func checkExtension(path string) {
	extension := filepath.Ext(path)
	if extension != ".kdbx" {
		log.Fatal("Unrecognized file extension!")
	}
}

func evaluateEntry(entry gokeepasslib.Entry, key string) bool {
	key = strings.ToLower(key)
	var username = entry.GetContent("UserName")
	var title = entry.GetTitle()

	var containsUsername = strings.Contains(strings.ToLower(username), key)
	var containsTitle = strings.Contains(strings.ToLower(title), key)

	return (containsUsername || containsTitle)

}

func findEntries(key string, groups []gokeepasslib.Group, path []string) []result {
	var entries []result

	for _, group := range groups {
		// enter the path and save it to the array
		path = append(path, group.Name)
		if len(group.Entries) > 0 {
			for _, entry := range group.Entries {
				if evaluateEntry(entry, key) {
					var localPath string = strings.Join(path[:], "/")
					entries = append(entries, result{entry: entry, parentGroup: localPath})
				}
			}
		}
		if len(group.Groups) > 0 {
			foundEntries := findEntries(key, group.Groups, path)
			if len(foundEntries) > 0 {
				for _, entry := range foundEntries {
					entries = append(entries, entry)
				}
			}
		}
		// exit the path
		path = path[:len(path)-1]
	}

	return entries
}

func clearClipboard(key string) {
	<-time.After(10 * time.Second)
	text, err := clipboard.ReadAll()
	if err != nil {
		return
	}
	if text == key {
		clipboard.WriteAll("")
	}
}

func copyUsername(entry gokeepasslib.Entry) {
	var username = entry.GetContent("UserName")

	if err := clipboard.WriteAll(username); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Username copied.")
}

func copyPassword(entry gokeepasslib.Entry) {
	var password = entry.GetPassword()

	if err := clipboard.WriteAll(password); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Password copied.")
	go clearClipboard(password)
}

func fileExists(file string) bool {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (data *keepass) IsOpen(file string) bool {
	file = os.ExpandEnv(file)
	for _, filepath := range data.Filepaths {
		if file == filepath {
			return true
		}
	}
	return false
}

func (data *keepass) findKey(key string) {
	var results []result
	data.lastResults = nil
	var path []string
	for _, DB := range data.DBs {
		results = findEntries(key, DB.Content.Root.Groups, path)
		data.lastResults = append(data.lastResults, results...)
	}
}

func (data *keepass) printSearchResults() {
	if len(data.lastResults) == 0 {
		fmt.Println("No results found.")
		return
	}
	fmt.Println("Search results:")
	for index, result := range data.lastResults {
		fmt.Printf("%d:  %s/%s\n", index, result.parentGroup, result.entry.GetTitle())
	}
}

func (data *keepass) readPassword() string {
	fmt.Print("Password (enter to abort): ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	return string(bytePassword)
}

func (data *keepass) storeDb(kdbxFilePath string) {
	var basename string = path.Base(kdbxFilePath)
	var dbName string = strings.TrimSuffix(basename, filepath.Ext(basename))
	data.DbNames = append(data.DbNames, dbName)
	data.Filepaths = append(data.Filepaths, kdbxFilePath)
}

// DBCount returns the number of DBs currently opened
func (data *keepass) DBCount() int {
	return len(data.DBs)
}

// Copy copies the username from the last results array.
// The given id is the index from the last results array.
func (data *keepass) Copy(id int, what string) {
	if len(data.lastResults) == 0 || (len(data.lastResults)-1) < id {
		fmt.Println("Invalid ID.")
		return
	}
	var result = data.lastResults[id].entry
	switch what {
	case "username":
		copyUsername(result)
	case "password":
		copyPassword(result)
	}
}

func removeElementFromStringSlice(id int, slice []string) []string {
	var newLen = len(slice) - 1
	copy(slice[id:], slice[id+1:])
	slice[newLen] = ""
	return slice[:newLen]
}

// Close closes the database with the given id
func (data *keepass) Close(id int) {
	var dbName = data.DbNames[id]

	data.Filepaths = removeElementFromStringSlice(id, data.Filepaths)
	data.DbNames = removeElementFromStringSlice(id, data.DbNames)

	newLen := len(data.DBs) - 1
	copy(data.DBs[id:], data.DBs[id+1:])
	data.DBs[newLen] = nil
	data.DBs = data.DBs[:newLen]
	data.lastResults = nil

	fmt.Printf("Database %s was successfully removed.\n", dbName)
}

// Open opens a keepass file
func (data *keepass) Open(filepath string) {
	checkExtension(filepath)

	filepath = os.ExpandEnv(filepath)
	if !fileExists(filepath) {
		fmt.Printf("\033[0;31m%s does not exist!\n\n\033[0m", filepath)
		return
	}
	if data.IsOpen(filepath) {
		fmt.Printf("\033[0;33mDatabase %s is already open!\n\033[0m", filepath)
		return
	}

	file, err := os.Open(filepath)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Opening file %s\n", filepath)

	var currentPassword = data.readPassword()

	if currentPassword == "" {
		fmt.Println("\033[0;33mFind the password and continue with opening the files by using the cont command.\033[0;0m")
		return
	}

	var currentDb = gokeepasslib.NewDatabase()
	currentDb.Credentials = gokeepasslib.NewPasswordCredentials(currentPassword)
	data.DBs = append(data.DBs, currentDb)

	decodeError := gokeepasslib.NewDecoder(file).Decode(currentDb)
	currentPassword = ""

	if decodeError != nil {
		fmt.Printf("\033[0;31m%s was not opened!\n\n\033[0m", filepath)
		return
	}

	data.storeDb(filepath)
	currentDb.UnlockProtectedEntries()

	fmt.Printf("\033[0;32mKdbx database %s opened successfully!\n\n\033[0m", filepath)
}

// Find finds the given key in the keepass file.
func (data *keepass) Find(key string) int {
	data.findKey(key)
	data.printSearchResults()
	return len(data.lastResults)
}

// PrintDatabases prints all opened databases to the prompt
func (data *keepass) PrintDatabases() {
	if len(data.DbNames) > 0 {
		fmt.Println("Currently opened KeePass databases:")
	} else {
		fmt.Println("No KeePass databases opened.")
		return
	}
	for index, dbName := range data.DbNames {
		fmt.Printf("%d: %s\t\t(%s)\n", index, dbName, data.Filepaths[index])
	}
}

// NewKeepassService is a provider for keepass.Service.
func NewKeepassService() interfaces.KeepassService {
	var keepassInstance = new(keepass)
	var service interfaces.KeepassService = keepassInstance
	return service
}
