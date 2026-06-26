package main

import (
	"fmt"
	"os"
	"os/exec"

	lib "github.com/apalala/jb/pkg"
)

const RealGitPath = "/opt/local/bin/git"

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
	args := os.Args[1:]
	for _, arg := range args {
		if len(arg) > 0 && arg[0] == '-' {
			continue
		}

		if ForbiddenCommands[arg] {
			lib.LogCmd()
			lib.Jb()
			os.Exit(0)
		}
		break
	}

	if err := lib.Call(RealGitPath, args...); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		fmt.Fprintln(os.Stderr, "Error executing Git:", err)
		os.Exit(1)
	}
}
