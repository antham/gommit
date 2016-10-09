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
	Use:   "check [commitFrom] [commitTo] [&path]",
	Short: "Check commit messages",
	Long: `check verify your commmits follow templates you defined
and return a list of commit that don't and exit with an error code.`,
	Run: func(cmd *cobra.Command, args []string) {
		from, to, path, err := extractArgs(args)

		if err != nil {
			failure(err)

			exitError()
		}

		infos, err := gommit.RunMatching(path, from, to, viper.GetStringMapString("matchers"))

		if err != nil {
			failure(err)

			exitError()
		}

		if len(*infos) != 0 {
			renderInfos(infos)
			renderExamples(viper.GetStringMapString("examples"))

			exitError()
		}

		success("Everything is ok")

		exitSuccess()
	},
}

func extractArgs(args []string) (string, string, string, error) {
	if len(args) < 2 {
		return "", "", "", fmt.Errorf("Two arguments required : origin commit and end commit")

	}

	if len(args) > 3 {
		return "", "", "", fmt.Errorf("3 arguments must be provided at most")
	}

	var path string

	if len(args) == 3 {
		f, err := os.Stat(args[2])

		if err != nil {
			return "", "", "", fmt.Errorf(`Ensure "%s" directory exists`, args[2])
		}

		if !f.IsDir() {
			return "", "", "", fmt.Errorf(`"%s" must be a directory`, args[2])
		}

		path = args[2]
	} else {
		var err error

		path, err = os.Getwd()

		if err != nil {
			return "", "", "", err
		}
	}

	return args[0], args[1], path, nil
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
