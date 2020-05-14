package main

import (
	"JudgerServer/container"
	"JudgerServer/router"
	"fmt"
	"os"

	seccomp "github.com/seccomp/libseccomp-golang"
)

func init() {
}

func main() {

    result := make(chan container.RunResult)
    go container.CreateRunner(&container.RunnerConfig{
        WorkDir:     "/home/darc/Code/JudgerServer/test_dir/",
        ChangeRoot:  true,
        GID:         1000,
        UID:         1000,
        RunablePath: "test",
        Arguments:   []string{
            // "g++",
            // "test2.cpp",
            // "-o",
            // "testqwq",
        },
        Envirment:       os.Environ(),
        OutputPath:      "output",
        InputPath:       "input",
        ErrorPath:       "error",
        SeccompRule:     container.DefaultSeccompBlacklist,
        SeccompType:     seccomp.ActAllow,
        RestrictExecve:  true,
        CPUTimeLimit:    5000,
        CompilerMode:    false,
        TimeLimit:       6000,
        MemoryLimit:     1024 * 1024 * 100,
        OutputSizeLimit: 1024,
    }, result)
    fmt.Println(<-result)
    router.Router.Run(":8081")
}
