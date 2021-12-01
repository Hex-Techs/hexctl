// Package run debug go process
package run

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/Hex-Techs/hexctl/pkg/common/file"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

const (
	MainFile = "main.go"
	WorkDir  = ".tmp"
)

var (
	cmd     *exec.Cmd
	history []int64
)

// Reload reload a go process
func Reload(command []string, startChan chan bool, stop chan int64) {
	<-startChan
	start(command)
	restart := true
	for mt := range stop {
		if len(history) == 2 {
			history[0], history[1] = history[1], mt
			ratelimit := history[1] - history[0]
			if ratelimit < 500 {
				restart = false
			} else {
				restart = true
			}
		}
		if len(history) == 0 || len(history) == 1 {
			history = append(history, mt)
			restart = true
		}
		if restart {
			color.Green.Println("Reload Progess...")
			restart = false
			Kill()
			time.Sleep(1 * time.Second)
			start(command)
		}
	}
}

// Kill kill a go process
func Kill() {
	pid := cmd.Process.Pid
	err := syscall.Kill(-pid, syscall.SIGINT)
	if err != nil {
		color.Red.Println("kill process with error:", err)
		if cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
			os.Exit(1)
		}
	}
	cmd.Process.Wait()
}

func start(command []string) {
	bin, err := build()
	cobra.CheckErr(err)
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
	color.Green.Println("Running process...")
	err = cmd.Start()
	cobra.CheckErr(err)
}

func build() (string, error) {
	color.Green.Println("Building...")
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
		color.Red.Println("building with error:", err)
		return "", nil
	}
	err = cmd.Wait()
	if err != nil {
		return "", err
	}
	hash2, err = file.Hash(bin)
	cobra.CheckErr(err)
	if hash == hash2 && hash != "" {
		color.Green.Println("No change found in project")
	}

	return bin, nil
}

// Clean clean the workdir
func Clean() {
	if err := os.RemoveAll(WorkDir); err != nil {
		color.Red.Println("clean with error:", err)
	}
}
