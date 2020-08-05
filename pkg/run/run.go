package run

import (
	"io"
	"os"
	"os/exec"

	"github.com/Fize/n/pkg/output"
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
			kill()
			start(command)
		}
	}
}

func kill() {
	err := cmd.Process.Kill()
	if err != nil {
		output.Errorln("kill process with error:", err)
	}
}

func start(command []string) {
	goCmd := []string{"run", "main.go"}
	goCmd = append(goCmd, command...)
	cmd = exec.Command("go", goCmd...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		output.Fatalln("run process with error:", err)
	}
}
