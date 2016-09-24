package cmd

import (
	"fmt"

	"github.com/antham/gommit/gommit"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "App version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(gommit.GetVersion())
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
