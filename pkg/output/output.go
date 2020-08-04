package output

import (
	"fmt"
	"os"

	"github.com/gookit/color"
)

func Infoln(msg ...interface{}) {
	fmt.Println(msg...)
}

func Infof(format string, msg ...interface{}) {
	fmt.Printf(format, msg...)
}

func Progressln(msg ...interface{}) {
	color.Yellow.Println(msg...)
}

func Progressf(format string, msg ...interface{}) {
	color.Yellow.Printf(format, msg...)
}

func Errorln(msg ...interface{}) {
	color.Red.Println(msg...)
}

func Errorf(format string, msg ...interface{}) {
	color.Red.Printf(format, msg...)
}

func Noteln(msg ...interface{}) {
	color.Cyan.Println(msg...)
}

func Notef(format string, msg ...interface{}) {
	color.Cyan.Printf(format, msg...)
}

func Successln(msg ...interface{}) {
	color.Green.Println(msg...)
}

func Successf(format string, msg ...interface{}) {
	color.Green.Printf(format, msg...)
}

func Fatalln(msg ...interface{}) {
	color.Red.Println(msg...)
	os.Exit(0)
}

func Fatalf(format string, msg ...interface{}) {
	color.Red.Printf(format, msg...)
	os.Exit(0)
}
