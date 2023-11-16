package display

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/manifoldco/promptui"
)

// the select ui max size
const maxSize = 30

// NewTerminalDisplay create a new terminal display
func NewTerminalDisplay(title string, size int, items ...string) *terminalDisplay {
	td := terminalDisplay{
		Title: title,
		Size:  size,
		Items: items,
	}
	if td.Size > maxSize {
		td.Size = maxSize
	}
	return &td
}

// TerminalDisplay terminal display
type terminalDisplay struct {
	// the title of the display
	Title string
	// when the display is a selector, the display size
	Size int
	// the items of the display, when the display is a selector, the items are the options
	Items []string
}

// Select terminal ui for selected
func (td *terminalDisplay) Select() string {
	searcher := func(input string, index int) bool {
		item := td.Items[index]
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(item, input)
	}
	prompt := promptui.Select{
		Size:     td.Size,
		Label:    td.Title,
		Items:    td.Items,
		Searcher: searcher,
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return result
}

// Confirm terminal ui for confirmed
func (td *terminalDisplay) Confirm() bool {
	prompt := promptui.Prompt{
		Label:     td.Title,
		IsConfirm: true,
		Default:   "N",
	}
	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
	}
	if strings.ToLower(result) == "y" {
		return true
	}
	return false
}

// Table terminal ui for table
func (td *terminalDisplay) Table(header table.Row, content ...[]interface{}) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(header)
	for _, v := range content {
		t.AppendRow(v)
	}
	t.Render()
}
