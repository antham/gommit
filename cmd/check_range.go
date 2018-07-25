package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/antham/gommit/gommit"
)

// checkRangeCmd represents the check command
var checkRangeCmd = &cobra.Command{
	Use:   "range [revisionfrom] [revisionTo] [&path]",
	Short: "Check messages in range",
	Run: func(cmd *cobra.Command, args []string) {
		err := validateFileConfig()

		if err != nil {
			failure(err)

			exitError()
		}

		from, to, path, err := extractCheckRangeArgs(args)

		if err != nil {
			failure(err)

			exitError()
		}

		q := gommit.RangeQuery{
			Path:     path,
			From:     from,
			To:       to,
			Matchers: viper.GetStringMapString("matchers"),
			Options:  buildOptions(),
		}

		matchings, err := gommit.MatchRangeQuery(q)

		processMatchResult(matchings, err, viper.GetStringMapString("examples"))
	},
}

func extractCheckRangeArgs(args []string) (string, string, string, error) {
	if len(args) < 2 {
		return "", "", "", fmt.Errorf("Two arguments required : origin commit and end commit")

	}

	if len(args) > 3 {
		return "", "", "", fmt.Errorf("3 arguments must be provided at most")
	}

	var path string
	var err error

	if len(args) == 3 {
		path = args[2]
	}

	path, err = parseDirectory(path)

	if err != nil {
		return "", "", "", err
	}

	return args[0], args[1], path, nil
}

func init() {
	checkCmd.AddCommand(checkRangeCmd)
}
