package safe

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/go-ps"
)

const LogFilePath = "/tmp/cmd.log"

func Jb() error {
  return Call("~/bin/jb")
}

func TriggerFence() {
	LogCmd()
	Jb()
	os.Exit(0)
}

func ExpandHome(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if path == "~" {
		return "/tmp"
	}
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:])
	}
	return path
}

func LogCmd() {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	args := os.Args
	var joinedArgs string
	if len(args) == 0 {
		joinedArgs = "[interactive shell or no args]"
	} else {
		joinedArgs = strings.Join(args, " ")
	}

	pid := os.Getpid()
	level := 0
	for pid != 0 {
		proc, err := ps.FindProcess(pid)
		if err != nil || proc == nil {
			break
		}

		cmdline := joinedArgs
		if level > 0 {
			cmdline = proc.Executable()
		}

		executalbe := fmt.Sprintf("(%s):", proc.Executable())
		logEntry := fmt.Sprintf("#%d [%s] PID %05d %16s %s\n", level, timestamp, pid, executalbe, cmdline)
		if f, err := os.OpenFile(LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			_, _ = f.WriteString(logEntry)
			f.Close()
		}

		pid = proc.PPid()
		level++
	}

}

func Call(cmd string, args ...string) error {
	exe := exec.Command(ExpandHome(cmd), args...)
	exe.Stdin = os.Stdin
	exe.Stdout = os.Stdout
	exe.Stderr = os.Stderr

	return exe.Run()
}

type SafeCfg struct {
	RealPath          string
	Name              string
	CmdFilter         map[string]bool
	OptFilter         map[string]bool
	ExtFilter         map[string]bool
	LogAll            bool
}

func SafeRun(cfg SafeCfg) {
	if cfg.LogAll {
		LogCmd()
	}
	good := true
	for _, arg := range os.Args[1:] {
		for suffix, bad := range cfg.ExtFilter {
			if bad && strings.HasSuffix(strings.ToLower(arg), suffix) {
				good = false
				break
			}
		}
		if !good {
			break
		}

		if len(arg) > 0 && arg[0] == '-' {
			for opt, bad := range cfg.OptFilter {
				if bad && strings.HasPrefix(arg, opt) {
					good = false
					break
				}
			}
			if !good {
				break
			}
			continue
		}

		if cfg.CmdFilter != nil && cfg.CmdFilter[arg] {
			good = false
			break
		}
	}
	if !good {
		TriggerFence()
	}

	if err := Call(cfg.RealPath, os.Args[1:]...); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		fmt.Fprintln(os.Stderr, "Error executing", cfg.Name+":", err)
		os.Exit(1)
	}
}
