package run

import (
	"bufio"
	"fmt"
	"io"
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
	f, ef := start(command)
	for {
		select {
		case <-stop:
			kill()
			f, ef = start(command)
		default:
			line, e := f.ReadString('\n')
			errline, ee := ef.ReadString('\n')
			switch e {
			case io.EOF:
				err := cmd.Wait()
				if err != nil {
					output.Errorf("an error occured while running process: %v\n", err)
				}
			case nil:
				fmt.Print(line)
			}
			switch ee {
			case io.EOF:
				err := cmd.Wait()
				if err != nil {
					output.Errorf("an error occured while running process: %v\n", err)
				}
			case nil:
				if errline != "" {
					fmt.Printf(errline)
				}
			}
		}
	}
}

func kill() {
	err := cmd.Process.Kill()
	if err != nil {
		output.Errorln("kill process with error:", err)
	}
}

func start(command []string) (*bufio.Reader, *bufio.Reader) {
	goCmd := []string{"run", "main.go"}
	goCmd = append(goCmd, command...)
	cmd = exec.Command("go", goCmd...)
	out, _ = cmd.StdoutPipe()
	errout, _ = cmd.StderrPipe()
	f := bufio.NewReader(out)
	ef := bufio.NewReader(errout)
	if err := cmd.Start(); err != nil {
		output.Fatalln("run process with error:", err)
	}
	return f, ef
}
