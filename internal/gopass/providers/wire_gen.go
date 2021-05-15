// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package providers

import (
	"gopass/gopass/internal/gopass/interfaces"
	"gopass/gopass/internal/gopass/services"
)

// Injectors from container.go:

func CreatePromptService(keepassService interfaces.KeepassService) interfaces.PromptService {
	multiFileService := services.NewMultifileService(keepassService)
	promptService := services.NewPromptService(keepassService, multiFileService)
	return promptService
}

func CreateKeepassService() interfaces.KeepassService {
	keepassService := services.NewKeepassService()
	return keepassService
}

func CreateMultifileService() interfaces.MultiFileService {
	keepassService := services.NewKeepassService()
	multiFileService := services.NewMultifileService(keepassService)
	return multiFileService
}

func CreateMultifileManager() *services.Manager {
	keepassService := services.NewKeepassService()
	multiFileService := services.NewMultifileService(keepassService)
	promptService := services.NewPromptService(keepassService, multiFileService)
	manager := services.NewMultifileManager(keepassService, promptService, multiFileService)
	return manager
}