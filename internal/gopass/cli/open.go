package cli

import (
	"github.com/spf13/cobra"
	"gopass/gopass/internal/gopass/providers"
)

var keepassService = providers.CreateKeepassService()
var kdbx string
var config string

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Opens a KeePass file",
	Long: `Opens a KeePass (*.kdbx) file and starts an
	interactive interface for working with the kdbx file.`,
	Run: func(cmd *cobra.Command, args []string) {
		keepassService.Open(kdbx)

		var promptService = providers.CreatePromptService(keepassService)
		promptService.Start()
	},
}

func init() {

	openCmd.Flags().StringVarP(&kdbx, "kdbx", "k", "", "kdbx file to open")
	openCmd.MarkFlagRequired("kdbx")

	rootCmd.AddCommand(openCmd)

}
