// Package run debug go process
package run

import (
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/Hex-Techs/hexctl/pkg/display"
)

var (
	cmd    *exec.Cmd
	out    io.Reader
	errout io.Reader
)

// Reload reload a go process
func Reload(command []string, stop chan bool) {
	start(command)
	for {
		select {
		case <-stop:
			Kill()
			time.Sleep(1 * time.Second)
			start(command)
		}
	}
}

// Kill kill a go process
func Kill() {
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		display.Errorln("kill process with error:", err)
	} else {
		syscall.Kill(-pgid, 15) // note the minus sign
	}
}

func start(command []string) {
	goCmd := []string{"run", "main.go"}
	goCmd = append(goCmd, command...)
	cmd = exec.Command("go", goCmd...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		display.Errorln("run process with error:", err)
	}
}
