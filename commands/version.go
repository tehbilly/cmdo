package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.0.1"

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Print cmdo's version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:", version)
	},
}

func init() {
	RootCmd.AddCommand(versionCommand)
}
