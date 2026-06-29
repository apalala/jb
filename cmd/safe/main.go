package main

import (
	"fmt"
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
		fmt.Fprintf(os.Stderr, "unknown: %s\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
}
