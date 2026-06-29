package main

import (
	lib "github.com/apalala/jb/pkg"
)

const RealPythonPath = "/usr/local/bin/python3"

func main() {
	// FIXME: add "-c" to forbiddenOptions when ready
	lib.SafeRun(lib.SafeCfg{
		RealPath:          RealPythonPath,
		Name:              "python3",
		ExtFilter: map[string]bool{".py": true},
	})
}
