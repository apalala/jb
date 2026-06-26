package main

import (
	"fmt"
	"os"
	"os/exec"

	lib "github.com/apalala/jb/pkg"
)

const RealBashPath = "/opt/local/bin/bash"

func main() {
	lib.LogCmd()

	args := os.Args[1:]
	if err := lib.Call(RealBashPath, args...); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		fmt.Fprintln(os.Stderr, "Error executing Bash:", err)
		os.Exit(1)
	}
}
