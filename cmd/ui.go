package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/antham/gommit/gommit"
)

var failure = func(err error) {
	color.Red(err.Error())
}

var success = func(message string) {
	color.Green(message)
}

var info = func(message string) {
	color.Cyan(message)
}

var renderMatchings = func(matchings *[]*gommit.Matching) {
	for _, m := range *matchings {
		color.White("----")
		fmt.Println()

		if ID, ok := m.Context["ID"]; ok {
			fmt.Printf("%s%s\n", color.YellowString("Id       : "), color.WhiteString("%s", ID))
		}

		if message, ok := m.Context["message"]; ok {
			color.Yellow("Message  : ")

			for _, field := range strings.Split(message, "\n") {
				fmt.Printf("%s%s\n", color.YellowString("           "), color.WhiteString("%s", field))
			}
		}

		fmt.Println()

		errs := []error{}

		if m.MessageError != nil {
			errs = append(errs, m.MessageError)
		}

		if m.SummaryError != nil {
			errs = append(errs, m.SummaryError)
		}

		for i, e := range errs {
			if i == 0 {
				fmt.Printf("%s", color.YellowString("Error(s) : "))
				fmt.Printf("- %s\n", color.RedString("%s", e.Error()))
			} else {
				fmt.Printf("           - %s\n", color.RedString("%s", e.Error()))
			}
		}

		fmt.Println()
	}
}

var renderExamples = func(examples map[string]string) {
	color.White("=======")
	fmt.Println()

	color.White("Your message must match one of those following patterns :")

	fmt.Println()

	for key, example := range examples {
		color.White("----")
		fmt.Println()

		color.Yellow("%s : ", strings.Replace(cases.Title(language.English).String(key), "_", " ", -1))
		fmt.Println()

		color.Cyan("%s", example)
	}
}
