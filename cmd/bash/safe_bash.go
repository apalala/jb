package main

import (
	lib "github.com/apalala/jb/pkg"
)

const RealBashPath = "/opt/local/bin/bash"

func main() {
	lib.SafeRun(lib.SafeCfg{
		RealPath: RealBashPath,
		Name:     "bash",
		LogAll:   true,
	})
}
