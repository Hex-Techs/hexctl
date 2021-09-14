package display

import (
	"strings"

	"github.com/gookit/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// Select terminal ui for selected
func Select(title string, size int, items []string) string {
	if size > 30 {
		size = 30
	}

	searcher := func(input string, index int) bool {
		item := items[index]
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(item, input)
	}

	prompt := promptui.Select{
		Size:     size,
		Label:    title,
		Items:    items,
		Searcher: searcher,
	}

	_, result, err := prompt.Run()
	checkErr(err)
	return result
}

// Confirm terminal ui for confirmed
func Confirm(title string) bool {
	prompt := promptui.Prompt{
		Label:     title,
		IsConfirm: true,
		Default:   "N",
	}

	result, err := prompt.Run()
	checkErr(err)
	choose := strings.ToLower(result)
	if choose == "y" || choose == "Y" {
		return true
	}
	return false
}

func checkErr(err error) {
	if err != nil {
		if err.Error() == "" {
			return
		}
		if err.Error() != "^C" {
			color.Red.Println(err)
			cobra.CheckErr(err)
		}
	}
}
