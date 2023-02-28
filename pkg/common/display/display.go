package display

import (
	"os"

	inf "github.com/fzdwx/infinite"
	"github.com/fzdwx/infinite/components"
	"github.com/fzdwx/infinite/components/input/confirm"
	"github.com/fzdwx/infinite/components/selection/singleselect"
	"github.com/jedib0t/go-pretty/v6/table"
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
	input := components.NewInput()
	_, err := inf.NewSingleSelect(td.Items, singleselect.WithFilterInput(input)).Display(td.Title)
	if err != nil {
		os.Exit(1)
	}
	return input.Value()
}

// Confirm terminal ui for confirmed
func (td *terminalDisplay) Confirm() bool {
	c := inf.NewConfirm(
		confirm.WithPure(),
		confirm.WithDefaultYes(),
		confirm.WithPrompt(td.Title),
		confirm.WithDisplayHelp(),
	)
	c.Display()
	return c.Value()
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
