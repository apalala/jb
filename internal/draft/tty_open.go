package draft

import (
	"os"
	"runtime"
)

func getTerminalDescriptor() string {
	if runtime.GOOS == "windows" {
		return "CONOUT$"
	}
	return "/dev/tty"
}

func ttyMain() {
	// sigChan := make(chan os.Signal, 1)

	// signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// go func() {
	// 	_ = <-sigChan
	// }()

	out := os.Stdout
    if stat, _ := out.Stat(); stat.Mode()&os.ModeCharDevice == 0 {
        termPath := getTerminalDescriptor()
        tty, err := os.OpenFile(termPath, os.O_WRONLY, 0)
        if err == nil {
            defer tty.Close()
            out = tty
        }
    }
}
