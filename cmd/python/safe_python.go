package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	lib "github.com/apalala/jb/pkg"
)

// Hardcoded path to the real, trusted Python interpreter
const RealPythonPath = "/Library/Frameworks/Python.framework/Versions/3.14/bin/python3"

func main() {
	// os.Args[0] is the program name; os.Args[1:] are the actual flags/scripts
	args := os.Args[1:]

	for _, arg := range args {
		// 1. Block inline command execution (-c, --command)
		if arg == "-c" || strings.HasPrefix(arg, "-c=") {
			lib.Jb()
			// fmt.Fprintln(os.Stderr, "Error: Inline Python execution (-c) is strictly forbidden.")
			os.Exit(0)
		}

		// 2. Block direct script execution (.py files)
		if strings.HasSuffix(strings.ToLower(arg), ".py") {
			lib.Jb()
			// fmt.Fprintf(os.Stderr, "Error: Direct .py script execution (%s) is forbidden.\n", arg)
			os.Exit(0)
		}
	}

	// If it passes validation, transparently hand off to the real Python interpreter
	cmd := exec.Command(RealPythonPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run and mirror the exit code
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		fmt.Fprintln(os.Stderr, "Error executing Python:", err)
		os.Exit(1)
	}
}
