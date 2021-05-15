package cli

import (
	"github.com/spf13/cobra"
	"gopass/gopass/internal/gopass/providers"
)

// manyCmd represents the openPreset command
var manyCmd = &cobra.Command{
	Use:   "many",
	Short: "Opens multiple kdbx files from a config file.",
	Long: `Allows reading multiple files at once and operate on them from a single CLI interface.
	Convenient if you are using multiple kdbx files on a regular basis.`,
	Run: func(cmd *cobra.Command, args []string) {
		var manager = providers.CreateMultifileManager()
		manager.Multifile.Open()
		manager.Prompt.Start()
	},
}

func init() {
	rootCmd.AddCommand(manyCmd)
}
