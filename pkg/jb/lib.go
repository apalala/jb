package jb

import (
	"os"
	"runtime"
	"syscall"

	"golang.org/x/sys/unix"
)

func preferredOut() *os.File {
	out := os.Stdout
	if isDevNull(out) {
		out = os.Stderr
	}
	if isDevNull(out) {
		termPath := getTerminalDescriptor()
		tty, err := os.OpenFile(termPath, os.O_WRONLY, 0)
		if err == nil {
			out = tty
		}
	}
	return out
}

func outputType(file *os.File) int {
	stat, _ := file.Stat()

	// Check if stdout is a standard interactive terminal
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return 0 // Interactive terminal
	} else if (stat.Mode() & os.ModeNamedPipe) != 0 {
		return 1 // Named pipe
	} else {
		return 2 // File or other
	}
}

func isDevNull(file *os.File) bool {
	stat, err := file.Stat()
	if err != nil {
		return false
	}

	sys, ok := stat.Sys().(*syscall.Stat_t)
	if !ok {
		return false
	}

	// Use the official x/sys/unix macros to parse device numbers
	major := unix.Major(uint64(sys.Rdev))
	minor := unix.Minor(uint64(sys.Rdev))

	switch runtime.GOOS {
	case "linux":
		return major == 1 && minor == 3
	case "darwin": // macOS
		return major == 3 && minor == 2
	default:
		return false
	}
}

func getTerminalDescriptor() string {
	if runtime.GOOS == "windows" {
		return "CONOUT$"
	}
	return "/dev/tty"
}
