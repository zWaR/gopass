package interfaces

// KeepassService is the interface implemented by the keepass.Service.
type KeepassService interface {
	Open(filepath string)
	Copy(id int, what string)
	Find(key string) int
	DBCount() int
	PrintDatabases()
	Close(id int)
	IsOpen(file string) bool
}

// PromptService is the interface implemented by prompt.Service.
type PromptService interface {
	Start()
}

// MultiFileService is the interface implemented by multifile.Service.
type MultiFileService interface {
	Open()
	InitConfig()
	LocationsCount() int
	ContinueOpen(file string)
}
