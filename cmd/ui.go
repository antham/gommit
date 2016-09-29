package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
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

func renderInfos(infos *[]map[string]string) {
	for _, info := range *infos {
		color.White("----")
		fmt.Println()

		fmt.Printf("%s%s\n", color.YellowString("Id      : "), color.WhiteString("%s", info["id"]))
		color.Yellow("Message : ")

		for _, field := range strings.Split(info["message"], "\n") {
			fmt.Printf("%s%s\n", color.YellowString("          "), color.WhiteString("%s", field))
		}
	}
}

func renderExamples(examples map[string]string) {
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
