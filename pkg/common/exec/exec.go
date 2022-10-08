package exec

import (
	"os"
	"os/exec"
)

// RunCommand run a shell command
func RunCommand(args string) error {
	cmd := exec.Command("sh", "-c", args)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
