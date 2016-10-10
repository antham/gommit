package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"

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

var renderErrors = func(errors *[]gommit.CommitError) {
	for _, e := range *errors {
		color.White("----")
		fmt.Println()

		fmt.Printf("%s%s\n", color.YellowString("Id       : "), color.WhiteString("%s", e.ID))

		color.Yellow("Message  : ")

		for _, field := range strings.Split(e.Message, "\n") {
			fmt.Printf("%s%s\n", color.YellowString("           "), color.WhiteString("%s", field))
		}

		fmt.Println()

		errs := []error{}

		if e.MessageError != nil {
			errs = append(errs, e.MessageError)
		}

		if e.SummaryError != nil {
			errs = append(errs, e.SummaryError)
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

	color.White("Your message commits must match one of those following patterns :")

	fmt.Println()

	for key, example := range examples {
		color.White("----")
		fmt.Println()

		color.Yellow("%s : ", strings.Replace(strings.Title(key), "_", " ", -1))
		fmt.Println()

		color.Cyan("%s", example)
	}
}
