package display

import (
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// SelectUI terminal ui for selected
func SelectUI(title string, items []string) string {
	prompt := promptui.Select{
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
	}

	result, err := prompt.Run()
	checkErr(err)
	choose := strings.ToLower(result)
	if choose == "y" || choose == "Y" || choose == "" {
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
			Errorln(err)
			cobra.CheckErr(err)
		}
	}
}
