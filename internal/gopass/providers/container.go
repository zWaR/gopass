//+build wireinject

package providers

import (
	"github.com/google/wire"
	"gopass/gopass/internal/gopass/interfaces"
	"gopass/gopass/internal/gopass/services"
)

func CreatePromptService(keepassService interfaces.KeepassService) interfaces.PromptService {
	panic(
		wire.Build(
			services.NewMultifileService,
			services.NewPromptService,
		),
	)
}

func CreateKeepassService() interfaces.KeepassService {
	panic(
		wire.Build(
			services.NewKeepassService,
		),
	)
}

func CreateMultifileService() interfaces.MultiFileService {
	panic(
		wire.Build(
			services.NewKeepassService,
			services.NewMultifileService,
		),
	)
}

func CreateMultifileManager() *services.Manager {
	panic(
		wire.Build(
			services.NewPromptService,
			services.NewKeepassService,
			services.NewMultifileService,
			services.NewMultifileManager,
		),
	)
}
