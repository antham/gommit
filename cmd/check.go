package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/antham/gommit/gommit"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check commit messages",
	Run: func(cmd *cobra.Command, args []string) {
		currentPath, err := os.Getwd()

		if err != nil {
			failure(err)

			exitError()
		}

		if len(args) != 2 {
			failure(fmt.Errorf("Two arguments required : origin commit and end commit"))

			exitError()
		}

		infos, err := gommit.RunMatching(currentPath, args[0], args[1], viper.GetStringMapString("matchers"))

		if err != nil {
			failure(err)

			exitError()
		}

		if len(*infos) != 0 {
			renderInfos(infos)
			renderExamples(viper.GetStringMapString("examples"))

			exitError()
		}

		success("Everyting ok")

		exitSuccess()
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
