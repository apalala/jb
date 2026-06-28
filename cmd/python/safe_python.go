package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	lib "github.com/apalala/jb/pkg"
)

const RealPythonPath = "/usr/local/bin/python3"

func main() {
	args := os.Args[1:]

	for _, arg := range args {
		// FIXME disable for now as unexpected programs depend on this
		// if arg == "-c" || strings.HasPrefix(arg, "-c=") {
		// 	lib.LogCmd()
		// 	lib.Jb()
		// 	os.Exit(0)
		// }

		if strings.HasSuffix(strings.ToLower(arg), ".py") {
			lib.LogCmd()
			lib.Jb()
			os.Exit(0)
		}
	}

	if err := lib.Call(RealPythonPath, args...); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		fmt.Fprintln(os.Stderr, "Error executing Python:", err)
		os.Exit(1)
	}
}
