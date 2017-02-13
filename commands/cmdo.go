package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "cmdo",
	Short: "A helpful all-in-one portable command thing!",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("I was called with args:")
		for _, arg := range args {
			fmt.Println("  -", arg)
		}

		fmt.Println("Available sub-commands:")
		for _, cmd := range cmd.Commands() {
			fmt.Println("  -", cmd.Name())
		}
	},
}
