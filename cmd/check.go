package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/antham/gommit/gommit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check ensure a message follows defined patterns",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()

		if err != nil {
			failure(err)

			exitError()
		}
	},
}

func parseDirectory(path string) (string, error) {
	if path == "" {
		if _, err := os.Getwd(); err != nil {
			return "", err
		}
	} else {
		f, err := os.Stat(path)

		if err != nil {
			return "", fmt.Errorf(`Ensure "%s" directory exists`, path)
		}

		if !f.IsDir() {
			return "", fmt.Errorf(`"%s" must be a directory`, path)
		}
	}

	return path, nil
}

func validateFileConfig() error {
	if len(viper.GetStringMapString("matchers")) == 0 {
		return fmt.Errorf("At least one matcher must be defined")
	}

	if len(viper.GetStringMapString("examples")) == 0 {
		return fmt.Errorf("At least one example must be defined")
	}

	for name, matcher := range viper.GetStringMapString("matchers") {
		_, err := regexp.Compile(matcher)

		if err != nil {
			return fmt.Errorf(`Regexp "%s" identified by "%s" is not a valid regexp, please check the syntax`, matcher, name)
		}
	}

	return nil
}

func processMatchResult(matchings *[]*gommit.Matching, err error, examples map[string]string) {
	if err != nil {
		failure(err)

		exitError()
	}

	if len(*matchings) != 0 {
		renderMatchings(matchings)
		renderExamples(examples)

		exitError()
	}

	success("Everything is ok")

	exitSuccess()
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
