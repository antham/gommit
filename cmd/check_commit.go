package cmd

import (
	"errors"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/antham/gommit/gommit"
)

// checkCommitCmd represents the command that check a commit message
var checkCommitCmd = &cobra.Command{
	Use:   "commit [id] [&path]",
	Short: "Check commit message",
	Run: func(cmd *cobra.Command, args []string) {
		err := validateFileConfig()
		if err != nil {
			failure(err)

			exitError()
		}

		ID, path, err := extractCheckCommitArgs(args)
		if err != nil {
			failure(err)

			exitError()
		}

		q := gommit.CommitQuery{
			ID:       ID,
			Path:     path,
			Matchers: viper.GetStringMapString("matchers"),
			Options:  buildOptions(),
		}

		matching, err := gommit.MatchCommitQuery(q)

		matchings := &[]*gommit.Matching{}

		if !gommit.IsZeroMatching(matching) {
			*matchings = append(*matchings, matching)
		}

		processMatchResult(matchings, err, viper.GetStringMapString("examples"))
	},
}

func extractCheckCommitArgs(args []string) (string, string, error) {
	if len(args) < 1 {
		return "", "", errors.New("one argument required : commit id")
	}

	ok, err := regexp.Match("[a-fA-F0-9]{40}", []byte(args[0]))

	if err != nil || !ok {
		return "", "", errors.New("argument must be a valid commit id")
	}

	if len(args) > 2 {
		return "", "", errors.New("2 arguments must be provided at most")
	}

	var path string

	if len(args) == 2 {
		path = args[1]
	}

	path, err = parseDirectory(path)
	if err != nil {
		return "", "", err
	}

	return args[0], path, nil
}

func init() {
	checkCmd.AddCommand(checkCommitCmd)
}
