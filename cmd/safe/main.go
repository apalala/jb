package main

import (
	"os"
	"path/filepath"

	"github.com/apalala/jb/pkg/safe"
)

func main() {
	switch filepath.Base(os.Args[0]) {
	case "git":
		safe.GitMain()
	case "python3":
		safe.PythonMain()
	case "bash":
		safe.BashMain()
	case "head":
		safe.HeadMain()
	default:
		safe.TriggerFence()
	}
}
