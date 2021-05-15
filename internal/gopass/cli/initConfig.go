package cli

import (
	"github.com/spf13/cobra"
	"gopass/gopass/internal/gopass/providers"
)

// writePresetCmd represents the writePreset command
var initConfigCmd = &cobra.Command{
	Use:   "initConfig",
	Short: "Initialize and create a config file",
	Long: `Allows reading multiple files at once and operate on them from a single CLI interface.
	Convenient if you are using multiple kdbx files on a regular basis.`,
	Run: func(cmd *cobra.Command, args []string) {
		var multifileService = providers.CreateMultifileService()
		multifileService.InitConfig()
	},
}

func init() {
	rootCmd.AddCommand(initConfigCmd)
}
