package utils

import (
	"io/ioutil"
	"os"

	"github.com/Hex-Techs/n/pkg/output"
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

func IsExists(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsDir(name string) bool {
	s, err := os.Stat(name)
	if err != nil {
		return false
	}
	return s.IsDir()
}
