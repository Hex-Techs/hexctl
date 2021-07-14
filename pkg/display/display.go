// Package display print the command line content
package display

import (
	"fmt"

	"github.com/gookit/color"
)

// Errorf print error message by given format
func Errorf(format string, msg ...interface{}) {
	color.Red.Printf(format, msg...)
}

// Errorln print error message by given message
func Errorln(msg ...interface{}) {
	color.Red.Println(msg...)
}

// Progressf print yellow message by given message
func Progressf(msg ...interface{}) {
	color.Yellow.Println(msg...)
}

// Progressln print yellow message by given message
func Progressln(msg ...interface{}) {
	color.Yellow.Println(msg...)
}

// Successln print sucessful message by given message
func Successln(msg ...interface{}) {
	color.Green.Println(msg...)
}

// Successf print sucessful message by given format
func Successf(format string, msg ...interface{}) {
	color.Green.Printf(format, msg...)
}

// Infoln print nomal message by given message
func Infoln(msg ...interface{}) {
	fmt.Println(msg...)
}

// Infof print nomal message by given format
func Infof(format string, msg ...interface{}) {
	fmt.Printf(format, msg...)
}
