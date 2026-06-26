package lib

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const LogFilePath = "/tmp/cmd.log"

func Jb() {
	cmd := exec.Command("/Users/apalala/bin/jb")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

func ExpandHome(path string) (string, error) {
	if path == "" {
		return "", nil
	}
	if path[0] != '~' {
		return path, nil
	}
	if len(path) == 1 || path[1] == '/' {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}

		return filepath.Join(home, path[1:]), nil
	}
	return path, nil
}

func LogCmd() {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	args := os.Args[1:]
	var joinedArgs string
	if len(args) == 0 {
		joinedArgs = "[interactive shell or no args]"
	} else {
		joinedArgs = strings.Join(args, " ")
	}

	logEntry := fmt.Sprintf("[%s] PID %d: %s\n", timestamp, os.Getpid(), joinedArgs)
	if f, err := os.OpenFile(LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		_, _ = f.WriteString(logEntry)
		f.Close()
	}

}

func main() {
	testPath := "~/projects/safe-bin"

	expanded, err := ExpandHome(testPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Original:", testPath)
	fmt.Println("Expanded:", expanded)
}
