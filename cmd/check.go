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
		err := validateConfig()

		if err != nil {
			failure(err)

			exitError()
		}

		from, to, path, err := extractArgs(args)

		if err != nil {
			failure(err)

			exitError()
		}

		errors, err := gommit.RunMatching(path, from, to, viper.GetStringMapString("matchers"), map[string]bool{
			"check-summary-length": viper.GetBool("config.check-summary-length"),
			"exclude-merge-commit": viper.GetBool("exclude-merge-commit"),
		})

		if err != nil {
			failure(err)

			exitError()
		}

		if len(*errors) != 0 {
			renderErrors(errors)
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

func validateConfig() error {
	if len(viper.GetStringMapString("matchers")) == 0 {
		return fmt.Errorf("At least one matcher must be defined")
	}

	if len(viper.GetStringMapString("examples")) == 0 {
		return fmt.Errorf("At least one example must be defined")
	}

	return nil
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
