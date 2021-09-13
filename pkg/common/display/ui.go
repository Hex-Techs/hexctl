package display

import (
	"strings"

	"github.com/gookit/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// SelectUI terminal ui for selected
func SelectUI(title string, size int, items []string) string {
	if size > 30 {
		size = 30
	}
	prompt := promptui.Select{
		Size:  size,
		Label: title,
		Items: items,
	}

	_, result, err := prompt.Run()
	checkErr(err)
	return result
}

// ConfirmUI terminal ui for confirmed
func ConfirmUI(title string) bool {
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
