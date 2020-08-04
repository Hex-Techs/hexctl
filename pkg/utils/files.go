package utils

import (
	"io/ioutil"
	"os"

	"github.com/Fize/n/pkg/output"
)

func Read(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		output.Errorln(err)
		os.Exit(1)
	}
	return string(content)
}

func Write(content, dst string) {
	if err := ioutil.WriteFile(dst, []byte(content), 0644); err != nil {
		output.Errorln(err)
		os.Exit(1)
	}
}
