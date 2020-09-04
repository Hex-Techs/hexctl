package utils

import (
	"os"
	"os/exec"

	"github.com/Fize/n/pkg/output"
)

func RunCommand(args string) {
	cmd := exec.Command("sh", "-c", args)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		output.Fatalln("run process with error:", err)
	}
	if err := cmd.Wait(); err != nil {
		output.Fatalln(err)
	}
}
