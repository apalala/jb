package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	lib "github.com/apalala/jb/pkg"
)

const (
	RealHeadPath   = "/usr/bin/head"
	ForbiddenMagic = "# Johannes Blues"
)

func main() {
	args := os.Args[1:]

	// 1. Check file arguments
	var hasFile bool
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			hasFile = true
			if err := inspectFile(arg); err != nil {
				// Let real head handle missing/unreadable files
			}
		}
	}

	// 2. Safely intercept stdin without deadlocking the pipeline
	if !hasFile {
		stat, err := os.Stdin.Stat()
		if err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
			inspectStdinAndPassThrough()
		}
	}

	// 3. Clear pass-through to the real binary
	if err := lib.Call(RealHeadPath, args...); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		fmt.Fprintln(os.Stderr, "Error executing head:", err)
		os.Exit(1)
	}
}

func inspectFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		if strings.Contains(scanner.Text(), ForbiddenMagic) {
			triggerFence()
		}
	}
	return nil
}

func inspectStdinAndPassThrough() {
	r, w, _ := os.Pipe()

	// Spawn a background worker to read from the old stdin, check it,
	// and stream it straight into the new stdin pipe without blocking.
	go func() {
		defer w.Close()

		reader := bufio.NewReader(os.Stdin)

		// Peek or read the first line safely
		firstLineBytes, err := reader.ReadBytes('\n')
		if len(firstLineBytes) > 0 {
			if strings.Contains(string(firstLineBytes), ForbiddenMagic) {
				triggerFence()
			}
			// Write the first line to our custom pipe right away
			w.Write(firstLineBytes)
		}

		if err == nil {
			// Stream the rest of the data continuously to avoid io.ReadAll deadlocks
			_, _ = io.Copy(w, reader)
		}
	}()

	// Replace original Stdin with the read end of our transparent proxy pipe
	os.Stdin = r
}

func triggerFence() {
	lib.LogCmd()
	lib.Jb()
	os.Exit(0)
}
