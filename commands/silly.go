package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(sillyOneCmd)
	RootCmd.AddCommand(sillyTwoCmd)
}

var sillyOneCmd = &cobra.Command{
	Use:   "silly-one",
	Short: "The first silly command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Silly one called!")
	},
}

var sillyTwoCmd = &cobra.Command{
	Use:   "silly-two",
	Short: "The first silly command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Silly two called!")
	},
}
