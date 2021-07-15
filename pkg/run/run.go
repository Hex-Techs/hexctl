// Package run debug go process
package run

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/Hex-Techs/hexctl/pkg/common/file"
	"github.com/Hex-Techs/hexctl/pkg/display"
	"github.com/spf13/cobra"
)

const (
	MainFile = "main.go"
	WorkDir  = ".tmp"
)

var (
	cmd *exec.Cmd
)

// Reload reload a go process
func Reload(command []string, stop chan bool) {
	start(command)
	for range stop {
		Kill()
		time.Sleep(1 * time.Second)
		start(command)
	}
}

// Kill kill a go process
func Kill() {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		display.Errorln("kill process with error:", err)
	} else {
		syscall.Kill(-pgid, 15) // note the minus sign
	}
}

func start(command []string) {
	bin, err := build()
	if err != nil {
		return
	}
	for {
		if file.IsExists(bin) {
			break
		}
	}
	if bin == "" {
		return
	}
	cmd = exec.Command(bin, command...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	display.Successln("Running process...")
	if err := cmd.Start(); err != nil {
		display.Errorln("run process with error:", err)
	}
}

func build() (string, error) {
	display.Successln("Building...")
	bin := fmt.Sprintf("%s/%s", WorkDir, "_main_")
	var hash, hash2 string
	var err error
	if file.IsExists(bin) {
		hash, err = file.Hash(bin)
		cobra.CheckErr(err)
	}

	cmd := exec.Command("go", []string{"build", "-o", bin}...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		display.Errorln("building with error:", err)
		return "", nil
	}
	err = cmd.Wait()
	if err != nil {
		return "", err
	}
	hash2, err = file.Hash(bin)
	cobra.CheckErr(err)
	if hash == hash2 && hash != "" {
		display.Successln("No change found in project")
	}

	return bin, nil
}
