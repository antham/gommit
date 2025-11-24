package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/antham/gommit/gommit"
)

// checkMessageCmd represents the command that check a message
var checkMessageCmd = &cobra.Command{
	Use:   "message [message]",
	Short: "Check message",
	Run: func(cmd *cobra.Command, args []string) {
		err := validateFileConfig()
		if err != nil {
			failure(err)

			exitError()
		}

		message, err := extractCheckMessageArgs(args)
		if err != nil {
			failure(err)

			exitError()
		}

		q := gommit.MessageQuery{
			Message:  message,
			Matchers: viper.GetStringMapString("matchers"),
			Options:  buildOptions(),
		}

		matching, err := gommit.MatchMessageQuery(q)

		matchings := &[]*gommit.Matching{}

		if !gommit.IsZeroMatching(matching) {
			*matchings = append(*matchings, matching)
		}

		processMatchResult(matchings, err, viper.GetStringMapString("examples"))
	},
}

func extractCheckMessageArgs(args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("one argument required : message")
	}

	return args[0], nil
}

func init() {
	checkCmd.AddCommand(checkMessageCmd)
}
