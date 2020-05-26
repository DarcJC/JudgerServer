package main

import (
	"JudgerServer/workspace"
	"fmt"
)

func init() {
}

func main() {

	// result := make(chan container.RunResult)
	// go container.CreateRunner(&container.RunnerConfig{
	//     WorkDir:     "/home/darc/Code/JudgerServer/test_dir/",
	//     ChangeRoot:  true,
	//     GID:         1000,
	//     UID:         1000,
	//     RunablePath: "test",
	//     Arguments:   []string{
	//         // "g++",
	//         // "test2.cpp",
	//         // "-o",
	//         // "testqwq",
	//     },
	//     Envirment:       os.Environ(),
	//     OutputPath:      "output",
	//     InputPath:       "input",
	//     ErrorPath:       "error",
	//     SeccompRule:     container.DefaultSeccompBlacklist,
	//     SeccompType:     seccomp.ActAllow,
	//     RestrictExecve:  true,
	//     CPUTimeLimit:    5000,
	//     CompilerMode:    false,
	//     TimeLimit:       6000,
	//     MemoryLimit:     1024 * 1024 * 100,
	//     OutputSizeLimit: 1024,
	// }, result)
	// fmt.Println(<-result)
	w := workspace.NewWorkSpace("/home/darc/Code/JudgerServer/test_dir", "", "", "", "cpp", "", false)
	w.MakeFiles()
	w.WriteSourceCode("JTIzaW5jbHVkZSUyMCUzQ2lvc3RyZWFtJTNFJTBBJTBBdXNpbmclMjBuYW1lc3BhY2UlMjBzdGQlM0IlMEElMEFpbnQlMjBtYWluJTI4JTI5JTIwJTdCJTBBJTIwY291dCUyMCUzQyUzQyUyMCUyMnF3cSUyMiUzQiUwQSUyMHJldHVybiUyMDAlM0IlMEElMEElN0Q=")
	fmt.Println(w.CompileSource())
	// router.Router.Run(":8081")
}
