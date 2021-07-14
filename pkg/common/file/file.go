package file

import (
	"crypto/sha256"
	"html/template"
	"io"
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
		return os.IsExist(err)
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

// WriteByTemp write a file by given template
func WriteByTemp(title, content, fileName string, data interface{}) {
	t := template.New(title)
	t = template.Must(t.Parse(content))
	var f *os.File
	var err error
	if !IsExists(fileName) {
		f, err = os.Create(fileName)
	} else {
		f, err = os.Open(fileName)
	}
	cobra.CheckErr(err)
	t.Execute(f, data)
}

func Hash(name string) (string, error) {
	file, err := os.Open(name)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return string(hash.Sum(nil)), nil
}
