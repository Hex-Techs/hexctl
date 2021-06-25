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

// Successln print sucessful message by given message
func Successln(msg ...interface{}) {
	color.Yellow.Println(msg...)
}

// Successf print sucessful message by given format
func Successf(format string, msg ...interface{}) {
	color.Yellow.Printf(format, msg...)
}

// Infoln print nomal message by given message
func Infoln(msg ...interface{}) {
	fmt.Println(msg...)
}

// Infof print nomal message by given format
func Infof(format string, msg ...interface{}) {
	fmt.Printf(format, msg...)
}
