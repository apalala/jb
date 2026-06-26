package main

import (
	"fmt"
	"os"
	"os/exec"

	lib "github.com/apalala/jb/pkg"
)

// Hardcoded path to the real, trusted Git binary
const RealGitPath = "/opt/local/bin/bash"

// A set of forbidden Git subcommands
var ForbiddenCommands = map[string]bool{
	"push":          true,
	"rebase":        true,
	"clone":         true,
	"remote":        true,
	"fetch":         true,
	"pull":          true,
	"submodule":     true,
	"filter-branch": true, // Prevents history rewriting
	"commit":        true,
	"checkout":      true,
	"co":            true,
}

func main() {
	lib.LogCmd()

	args := os.Args[1:]
	for _, arg := range args {
		if len(arg) > 0 && arg[0] == '-' {
			continue
		}

		if ForbiddenCommands[arg] {
			cmd := exec.Command("/Users/apalala/bin/jb")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			_ = cmd.Run()
			os.Exit(0)
		}
		break
	}

	cmd := exec.Command(RealGitPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		fmt.Fprintln(os.Stderr, "Error executing Git:", err)
		os.Exit(1)
	}
}
