package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
)

/*
#include <unistd.h>
#include <sched.h>

*/
import "C"

// RunnerConfig CreateRunner的配置项
type RunnerConfig struct {
	WorkDir    string
	ChangeRoot bool
	GID        int
	UID        int
}

// MapUser 映射命名空间内外的用户
func MapUser(uid, gid, pid int) {
	ufile := fmt.Sprintf("/proc/%d/uid_map", pid)
	udata := []byte(fmt.Sprintf("%d %d %d", 1, uid, 1))
	if err := ioutil.WriteFile(ufile, udata, 0666); err != nil {
		panic(err)
	}

	gfile := fmt.Sprintf("/proc/%d/gid_map", pid)
	gdata := []byte(fmt.Sprintf("%d %d %d", 1, gid, 1))
	if err := ioutil.WriteFile(gfile, gdata, 0666); err != nil {
		panic(err)
	}
}

// CreateRunner 创建运行进程
func CreateRunner(config *RunnerConfig) {
	// 初始化通讯管道
	pipefd := make([]int, 2)
	syscall.Pipe(pipefd)

	if err := syscall.Unshare(syscall.CLONE_NEWPID); err != nil {
		panic(err)
	}

	pid := Fork()
	if pid == 0 {
		// 子进程
		// 等待父进程
		if err := syscall.Close(pipefd[1]); err != nil {
			panic(err)
		}
		if _, err := syscall.Read(pipefd[0], []byte{1}); err != nil {
			panic(err)
		}

		// 切换工作目录
		if err := syscall.Chdir(config.WorkDir); err != nil {
			panic(err)
		}

		// 隔离命名空间
		if err := syscall.Unshare(syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER); err != nil {
			panic(err)
		}

		// 创建文件夹
		if err := os.MkdirAll(config.WorkDir+"/proc", 0666); err != nil {
			panic(err)
		}
		if err := os.MkdirAll(config.WorkDir+"/dev", 0666); err != nil {
			panic(err)
		}
		if err := os.MkdirAll(config.WorkDir+"/bin", 0666); err != nil {
			panic(err)
		}
		if err := os.MkdirAll(config.WorkDir+"/lib", 0666); err != nil {
			panic(err)
		}
		if err := os.MkdirAll(config.WorkDir+"/usr/lib", 0666); err != nil {
			panic(err)
		}
		if err := os.MkdirAll(config.WorkDir+"/var/lib", 0666); err != nil {
			panic(err)
		}

		// 重新挂载部分文件系统
		if err := syscall.Mount("proc", "/proc", "proc", syscall.MS_PRIVATE, ""); err != nil {
			panic(err)
		}
		if err := syscall.Mount("udev", "/dev", "devtmpfs", syscall.MS_PRIVATE, ""); err != nil {
			panic(err)
		}

		// chroot jail
		if config.ChangeRoot {
			// 绑定挂载部分文件夹
			if err := syscall.Mount("/usr/lib", config.WorkDir+"/usr/lib", "none", syscall.MS_BIND, ""); err != nil {
				panic(err)
			}
			if err := syscall.Mount("/lib", config.WorkDir+"/lib", "none", syscall.MS_BIND, ""); err != nil {
				panic(err)
			}
			if err := syscall.Mount("/bin", config.WorkDir+"/bin", "none", syscall.MS_BIND, ""); err != nil {
				panic(err)
			}

			if err := syscall.Chroot("./"); err != nil {
				panic(err)
			}
		}

		// 重定向IO流
		inputfd, err := syscall.Open("./input", syscall.O_RDONLY, 0666)
		if err != nil {
			panic(err)
		}
		outputfd, err := syscall.Open("./output", syscall.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		errorfd, err := syscall.Open("./error", syscall.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		if err := syscall.Dup2(inputfd, int(os.Stdin.Fd())); err != nil {
			panic(err)
		}
		if err := syscall.Dup2(outputfd, int(os.Stdout.Fd())); err != nil {
			panic(err)
		}
		if err := syscall.Dup2(errorfd, int(os.Stderr.Fd())); err != nil {
			panic(err)
		}

		// EXECVE子进程
		if err := syscall.Exec(config.WorkDir+"/test", []string{}, os.Environ()); err != nil {
			panic(err)
		}

		// 子进程execve失败 退出
		os.Exit(-1)
	} else if pid > 0 {
		// 父进程
		// MapUser(config.UID, config.GID, pid)
		// 通知子进程
		if err := syscall.Close(pipefd[1]); err != nil {
			panic(err)
		}

		// 等待子进程
		var wstatus *syscall.WaitStatus
		rusage := syscall.Rusage{}
		wpid, err := syscall.Wait4(pid, wstatus, 0, &rusage)

		if err != nil {
			panic(err)
		}

		fmt.Println("子进程退出：", wpid)

		return
	} else {
		// 运行错误
		panic("fork failed")
	}
}
