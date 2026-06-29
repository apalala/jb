package safe

const RealPythonPath = "/usr/local/bin/python3"

func PythonMain() {
	SafeRun(SafeCfg{
		RealPath:  RealPythonPath,
		Name:      "python3",
		ExtFilter: map[string]bool{".py": true},
	})
}
