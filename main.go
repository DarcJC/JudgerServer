package main

import (
	"JudgerServer/container"
	"JudgerServer/router"
	"os"

	seccomp "github.com/seccomp/libseccomp-golang"
)

func init() {
}

func main() {

	go container.CreateRunner(&container.RunnerConfig{
		WorkDir:     "/home/darc/Code/JudgerServer/test_dir/",
		ChangeRoot:  true,
		GID:         1000,
		UID:         1000,
		RunablePath: "/usr/bin/g++",
		Arguments: []string{
			"g++",
			"test.cpp",
			"-o",
			"testqwq",
		},
		Envirment:      os.Environ(),
		OutputPath:     "output",
		InputPath:      "input",
		ErrorPath:      "error",
		SeccompRule:    container.DefaultSeccompBlacklist,
		SeccompType:    seccomp.ActAllow,
		RestrictExecve: false,
		CPUTimeLimit:   30000,
		CompilerMode:   true,
	})
	router.Router.Run(":8081")
}
