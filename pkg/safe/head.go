package safe

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const (
	RealHeadPath   = "/usr/bin/head"
	ForbiddenMagic = "# Johannes Blues"
)

func HeadMain() {
	args := os.Args[1:]

	var hasFile bool
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			hasFile = true
			if err := inspectFile(arg); err != nil {
			}
		}
	}

	if !hasFile {
		stat, err := os.Stdin.Stat()
		if err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
			inspectStdinAndPassThrough()
		}
	}

	if err := Call(RealHeadPath, args...); err != nil {
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
			TriggerFence()
		}
	}
	return nil
}

func inspectStdinAndPassThrough() {
	r, w, _ := os.Pipe()

	go func() {
		defer w.Close()

		reader := bufio.NewReader(os.Stdin)

		firstLineBytes, err := reader.ReadBytes('\n')
		if len(firstLineBytes) > 0 {
			if strings.Contains(string(firstLineBytes), ForbiddenMagic) {
				TriggerFence()
			}
			w.Write(firstLineBytes)
		}

		if err == nil {
			_, _ = io.Copy(w, reader)
		}
	}()

	os.Stdin = r
}
