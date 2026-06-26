package main

import (
	"fmt"
	"os"
	"os/exec"

	lib "github.com/apalala/jb/pkg"
)

// Hardcoded path to the real, trusted system Bash binary
const RealBashPath = "/bin/bash"

func main() {
	lib.LogCmd()

	args := os.Args[1:]
	cmd := exec.Command(RealBashPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		fmt.Fprintln(os.Stderr, "Error executing Bash:", err)
		os.Exit(1)
	}
}
