package file

import (
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// Read read a file content form filepath
func Read(filePath string) string {
	content, err := ioutil.ReadFile(filePath)
	cobra.CheckErr(err)
	return string(content)
}

// Write write a file by given content
func Write(content, dst string) {
	err := ioutil.WriteFile(dst, []byte(content), 0644)
	cobra.CheckErr(err)
}

// IsExists the given file is exist
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

// IsDir the given is a directory or not
func IsDir(name string) bool {
	s, err := os.Stat(name)
	if err != nil {
		return false
	}
	return s.IsDir()
}
