package workspace

import (
	"JudgerServer/container"
	"fmt"
	"os"

	seccomp "github.com/seccomp/libseccomp-golang"
)

// CppCompiler C++编译
type CppCompiler struct {
	AbstractCompiler
}

// Compile 编译
func (c *CppCompiler) Compile(a *interface{}) (*container.RunResult, error) {
	args := []string{
		"g++",
		"-std=c++14",
		"-O2",
		"-o",
		"binary",
		"-x",
		"c++",
		c.Src,
	}

	result := make(chan container.RunResult)
	go container.CreateRunner(&container.RunnerConfig{
		WorkDir:         c.WS.BaseDir,
		ChangeRoot:      true,
		GID:             1000,
		UID:             1000,
		RunablePath:     "/usr/bin/g++",
		Arguments:       args,
		Envirment:       os.Environ(),
		OutputPath:      "output.compile",
		InputPath:       "input.compile",
		ErrorPath:       "error.compile",
		SeccompRule:     container.DefaultSeccompBlacklist,
		SeccompType:     seccomp.ActAllow,
		RestrictExecve:  false,
		CPUTimeLimit:    20000,
		CompilerMode:    true,
		TimeLimit:       20000,
		MemoryLimit:     1024 * 1024 * 100,
		OutputSizeLimit: 0,
	}, result)
	r := <-result
	fmt.Println(r)
	return &r, nil
}
