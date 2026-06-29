package safe

const RealBashPath = "/opt/local/bin/bash"

func BashMain() {
	SafeRun(SafeCfg{
		RealPath: RealBashPath,
		Name:     "bash",
		LogAll:   true,
	})
}
