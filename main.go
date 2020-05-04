package main

import (
	"JudgerServer/container"
)

func init() {
}

func main() {
	container.CreateRunner(&container.RunnerConfig{
		WorkDir:    "/home/darc/Code/JudgerServer/test_dir",
		ChangeRoot: false,
		GID:        1000,
		UID:        1000,
	})
	// router.Router.Run()
}
