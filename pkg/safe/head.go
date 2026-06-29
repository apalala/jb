package safe

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/apalala/jb/pkg/jb"
)

const (
	RealHeadPath = "/usr/bin/head"
)

func HeadMain() {
	args := os.Args[1:]

	hasFile := false
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			hasFile = true
			break
		}
	}

	if !hasFile {
		stat, err := os.Stdin.Stat()
		if err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
			if hasHeader() {
				TriggerFence()
				return
			}
		}
	}

	if err := Call(RealHeadPath, args...); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		os.Exit(1)
	}
}

func hasHeader() bool {
	reader := bufio.NewReader(os.Stdin)
	firstLineBytes, err := reader.ReadBytes('\n')
	if err != nil && len(firstLineBytes) == 0 {
		return false
	}

	found := strings.Contains(string(firstLineBytes), jb.JbHeader)

	r, w, _ := os.Pipe()
	go func() {
		defer w.Close()
		w.Write(firstLineBytes)
		if err == nil {
			_, _ = io.Copy(w, reader)
		}
	}()
	os.Stdin = r

	return found
}
