package utils

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// RunCommand run a shell command
func RunCommand(args string) {
	cmd := exec.Command("sh", "-c", args)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	cobra.CheckErr(err)
	err = cmd.Wait()
	cobra.CheckErr(err)
}
