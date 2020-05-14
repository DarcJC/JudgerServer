package container

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"syscall"
	"time"
	"unsafe"

	seccomp "github.com/seccomp/libseccomp-golang"
)

/*
#include <unistd.h>
#include <sched.h>
#include <stdio.h>

void run_child(char *path, char *args, char *envs) {
    execve(path, args, envs);
}

*/
import "C"

// RunnerConfig CreateRunner的配置项
type RunnerConfig struct {
	WorkDir        string // 必填 工作目录
	ChangeRoot     bool
	GID            int                   // 暂时无效
	UID            int                   // 暂时无效
	Arguments      []string              // 运行参数
	Envirment      []string              // 运行环境变量
	RunablePath    string                // 可执行文件路径, 基于WorkDir填写
	OutputPath     string                // 必填 输出文件路径
	InputPath      string                // 必填 输入文件路径
	ErrorPath      string                // 必填 错误输出文件路径
	SeccompRule    []seccomp.ScmpSyscall // 必填
	SeccompType    seccomp.ScmpAction    // seccomp.ActKill 或 seccomp.ActAllow 其它暂不支持
	RestrictExecve bool                  // 是否限制execve路径 g++会调用execve...
	CompilerMode   bool                  // 是否为编译器模式(不会进入新的命名空间)
	// 资源限制
	// <= 0 则 UNLIMITED
	MemoryLimit        int64  // Byte 内存限制
	CPUTimeLimit       int64  // ms CPU时间限制
	ProcessNumberLimit int64  // 个 最大创建进程数限制
	OutputSizeLimit    int64  // Byte 最大的输出文件大小
	CoreDumpLimit      uint64 // Byte 最大的核心转储大小 为0则禁用
	StackLimit         int64  // Byte 栈大小限制
	TimeLimit          int64  // ms 实际时间限制
	TimeAccuracy       uint64 // ms 计时器时间精度 毫秒
}

// RunResult 运行结果
type RunResult struct {
	WorkDir    string
	CPUTime    uint64 // ms
	RealTime   uint64 // ms
	Memory     uint64 // byte
	OutputPath string
	ErrorPath  string
	ExitCode   int
	Signal     int
	Result     string // 退出原因 不存在则为signal -1
}

// DefaultSeccompBlacklist 默认黑名单
var DefaultSeccompBlacklist []seccomp.ScmpSyscall

func init() {
	DefaultSeccompBlacklist = []seccomp.ScmpSyscall{
		GetSyscallNumber("acct"),
		GetSyscallNumber("add_key"),
		GetSyscallNumber("bpf"),
		GetSyscallNumber("clock_adjtime"),
		GetSyscallNumber("clock_settime"),
		// GetSyscallNumber("clone"),
		GetSyscallNumber("chroot"),
		// GetSyscallNumber("chdir"),
		GetSyscallNumber("create_module"),
		GetSyscallNumber("delete_module"),
		GetSyscallNumber("execveat"),
		GetSyscallNumber("finit_module"),
		GetSyscallNumber("get_kernel_syms"),
		GetSyscallNumber("get_mempolicy"),
		GetSyscallNumber("init_module"),
		GetSyscallNumber("ioperm"),
		GetSyscallNumber("iopl"),
		GetSyscallNumber("kcmp"),
		GetSyscallNumber("kexec_file_load"),
		GetSyscallNumber("kexec_load"),
		GetSyscallNumber("keyctl"),
		GetSyscallNumber("lookup_dcookie"),
		GetSyscallNumber("mbind"),
		GetSyscallNumber("mount"),
		GetSyscallNumber("move_pages"),
		GetSyscallNumber("name_to_handle_at"),
		GetSyscallNumber("nfsservctl"),
		GetSyscallNumber("open_by_handle_at"),
		GetSyscallNumber("perf_event_open"),
		GetSyscallNumber("personality"),
		GetSyscallNumber("pivot_root"),
		GetSyscallNumber("process_vm_readv"),
		GetSyscallNumber("process_vm_writev"),
		GetSyscallNumber("ptrace"),
		GetSyscallNumber("query_module"),
		GetSyscallNumber("quotactl"),
		GetSyscallNumber("reboot"),
		GetSyscallNumber("request_key"),
		GetSyscallNumber("set_mempolicy"),
		GetSyscallNumber("setns"),
		GetSyscallNumber("settimeofday"),
		GetSyscallNumber("setrlimit"),
		GetSyscallNumber("stime"),
		GetSyscallNumber("swapon"),
		GetSyscallNumber("swapoff"),
		GetSyscallNumber("sysfs"),
		GetSyscallNumber("_sysctl"),
		GetSyscallNumber("umount"),
		GetSyscallNumber("umount2"),
		GetSyscallNumber("unshare"),
		GetSyscallNumber("uselib"),
		GetSyscallNumber("userfaultfd"),
		GetSyscallNumber("ustat"),
		GetSyscallNumber("vm86"),
		GetSyscallNumber("vm86old"),
	}
}

// MapUser 映射命名空间内外的用户
// 不知道为啥会Operation not permitted
func MapUser(uid, gid, pid int) {
	ufile := fmt.Sprintf("/proc/%d/uid_map", pid)
	udata := []byte(fmt.Sprintf("%d %d %d", 0, uid, 0))
	if err := ioutil.WriteFile(ufile, udata, 0777); err != nil {
		panic(err)
	}

	gfile := fmt.Sprintf("/proc/%d/gid_map", pid)
	gdata := []byte(fmt.Sprintf("%d %d %d", 0, gid, 0))
	if err := ioutil.WriteFile(gfile, gdata, 0777); err != nil {
		panic(err)
	}
}

// If 假装有三目运算符
func If(b bool, t, f interface{}) interface{} {
	if b {
		return t
	}
	return f
}

// GetDefaultAccuracy 获取默认时钟精度
func GetDefaultAccuracy() uint64 {
	return 100
}

// NewWatcher 创建一个监控例程
func NewWatcher(pid int, config *RunnerConfig, quit <-chan int) {
	var counter uint64
	counter = 0
	accu := config.TimeAccuracy
	if accu <= 0 {
		accu = GetDefaultAccuracy()
	}
	ticker := time.NewTicker(time.Millisecond * time.Duration(accu))
	for {
		select {
		case <-quit:
			return // 收到退出信号
		case <-ticker.C:
			counter++
			if int64(counter*accu) > config.TimeLimit {
				if err := syscall.Kill(pid, 9); err != nil {
					panic(err)
				}
				return
			}
		}
	}
}

// CreateRunner 创建运行进程
func CreateRunner(config *RunnerConfig, res chan RunResult) {
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

		if !config.CompilerMode {
			// 隔离命名空间
			if err := syscall.Unshare(syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER | syscall.CLONE_FILES); err != nil {
				panic(err)
			}

			// 创建文件夹
			if err := os.MkdirAll(config.WorkDir+"/proc", 0777); err != nil {
				panic(err)
			}
			if err := os.MkdirAll(config.WorkDir+"/dev", 0777); err != nil {
				panic(err)
			}
			if err := os.MkdirAll(config.WorkDir+"/bin", 0777); err != nil {
				panic(err)
			}
			if err := os.MkdirAll(config.WorkDir+"/lib", 0777); err != nil {
				panic(err)
			}
			if err := os.MkdirAll(config.WorkDir+"/lib64", 0777); err != nil {
				panic(err)
			}
			if err := os.MkdirAll(config.WorkDir+"/usr/lib", 0777); err != nil {
				panic(err)
			}
			if err := os.MkdirAll(config.WorkDir+"/usr/bin", 0777); err != nil {
				panic(err)
			}
			if err := os.MkdirAll(config.WorkDir+"/usr/include", 0777); err != nil {
				panic(err)
			}
			if err := os.MkdirAll(config.WorkDir+"/usr/local/include", 0777); err != nil {
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
				if err := syscall.Mount("/lib64", config.WorkDir+"/lib64", "none", syscall.MS_BIND, ""); err != nil {
					panic(err)
				}
				if err := syscall.Mount("/bin", config.WorkDir+"/bin", "none", syscall.MS_BIND, ""); err != nil {
					panic(err)
				}
				if err := syscall.Mount("/usr/bin", config.WorkDir+"/usr/bin", "none", syscall.MS_BIND, ""); err != nil {
					panic(err)
				}
				if err := syscall.Mount("/usr/include", config.WorkDir+"/usr/include", "none", syscall.MS_BIND, ""); err != nil {
					panic(err)
				}
				if err := syscall.Mount("/usr/local/include", config.WorkDir+"/usr/local/include", "none", syscall.MS_BIND, ""); err != nil {
					panic(err)
				}

				if err := syscall.Chroot("./"); err != nil {
					panic(err)
				}
			}
		}

		// 重定向IO流
		inputfd, err := syscall.Open(config.InputPath, syscall.O_RDONLY, 0666)
		if err != nil {
			panic(err)
		}
		outputfd, err := syscall.Open(config.OutputPath, syscall.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
		errorfd, err := syscall.Open(config.ErrorPath, syscall.O_WRONLY, 0666)
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

		// 设置资源限制
		if config.MemoryLimit > 0 {
			if err := syscall.Setrlimit(syscall.RLIMIT_AS, &syscall.Rlimit{
				Cur: uint64(config.MemoryLimit * 2),
				Max: uint64(config.MemoryLimit * 2),
			}); err != nil {
				panic(err)
			}
		}
		if config.CPUTimeLimit > 0 {
			// CPU时间额外给出1秒
			if err := syscall.Setrlimit(syscall.RLIMIT_CPU, &syscall.Rlimit{
				Cur: uint64(config.CPUTimeLimit+1000) / 1000,
				Max: uint64(config.CPUTimeLimit+1000) / 1000,
			}); err != nil {
				panic(err)
			}
		}
		if config.OutputSizeLimit > 0 {
			if err := syscall.Setrlimit(syscall.RLIMIT_FSIZE, &syscall.Rlimit{
				Cur: uint64(config.OutputSizeLimit),
				Max: uint64(config.OutputSizeLimit),
			}); err != nil {
				panic(err)
			}
		}
		if config.StackLimit > 0 {
			if err := syscall.Setrlimit(syscall.RLIMIT_STACK, &syscall.Rlimit{
				Cur: uint64(config.StackLimit),
				Max: uint64(config.StackLimit),
			}); err != nil {
				panic(err)
			}
		}
		if err := syscall.Setrlimit(syscall.RLIMIT_CORE, &syscall.Rlimit{
			Cur: uint64(config.CoreDumpLimit),
			Max: uint64(config.CoreDumpLimit),
		}); err != nil {
			panic(err)
		}

		// Seccomp 规则
		filter, err := seccomp.NewFilter(config.SeccompType)
		if err != nil {
			panic(err)
		}
		if config.SeccompType == seccomp.ActAllow && config.SeccompRule == nil {
			config.SeccompRule = DefaultSeccompBlacklist
		}
		for _, s := range config.SeccompRule {
			if config.SeccompType == seccomp.ActAllow {
				filter.AddRule(s, seccomp.ActKill)
			} else {
				filter.AddRule(s, seccomp.ActAllow)
			}
		}
		targetPath := C.CString(config.RunablePath)
		if config.RestrictExecve {
			if config.SeccompType == seccomp.ActKill {
				execveAllow, err := seccomp.MakeCondition(0, seccomp.CompareEqual, uint64((uintptr)(unsafe.Pointer(targetPath))))
				if err != nil {
					panic(err)
				}
				if err := filter.AddRuleConditional(GetSyscallNumber("execve"), seccomp.ActAllow, []seccomp.ScmpCondition{execveAllow}); err != nil {
					panic(err)
				}
			} else {
				execveDeny, err := seccomp.MakeCondition(0, seccomp.CompareNotEqual, uint64((uintptr)(unsafe.Pointer(targetPath))))
				if err != nil {
					panic(err)
				}
				if err := filter.AddRuleConditional(GetSyscallNumber("execve"), seccomp.ActKill, []seccomp.ScmpCondition{execveDeny}); err != nil {
					panic(err)
				}
			}
		}
		if err := filter.Load(); err != nil {
			filter.Release()
			panic(err)
		}
		filter.Release()

		// EXECVE子进程
		args, err := syscall.SlicePtrFromStrings(config.Arguments)
		if err != nil {
			panic(err)
		}
		envs, err := syscall.SlicePtrFromStrings(config.Envirment)
		if err != nil {
			panic(err)
		}
		C.run_child(targetPath, (*C.char)(unsafe.Pointer(&args[0])), (*C.char)(unsafe.Pointer(&envs[0])))

		// 子进程execve失败 退出
		os.Exit(2133)
	} else if pid > 0 {
		// 父进程
		// MapUser(config.UID, config.GID, pid)
		// 通知子进程
		if err := syscall.Close(pipefd[1]); err != nil {
			panic(err)
		}

		// 开始时间
		start := time.Now()

		// 监控进程
		quit := make(chan int)
		if config.TimeLimit > 0 {
			go NewWatcher(pid, config, quit)
		}

		// 等待子进程
		wstatus := new(syscall.WaitStatus)
		rusage := syscall.Rusage{}
		wpid, _ := syscall.Wait4(pid, wstatus, 0, &rusage)
		close(quit) // 通知监控进程退出

		if res == nil {
			log.Printf("子进程退出：%d 状态：%d 信号: %d %s\n", wpid, wstatus.ExitStatus(), wstatus.Signal(), wstatus.Signal().String())
		} else {
			result := RunResult{}
			result.WorkDir = config.WorkDir
			result.OutputPath = config.OutputPath
			result.ErrorPath = config.ErrorPath
			result.Memory = uint64(rusage.Maxrss * 1000)
			result.CPUTime = uint64((rusage.Utime.Sec * 1000) + (rusage.Utime.Usec / 1000) + (rusage.Stime.Sec * 1000) + (rusage.Stime.Usec / 1000))
			result.RealTime = uint64(time.Since(start).Milliseconds())
			result.Signal = int(wstatus.Signal())
			result.ExitCode = wstatus.ExitStatus()
			result.Result = wstatus.Signal().String()
			res <- result
		}

		return
	} else {
		// Fork错误
		panic("fork failed")
	}
}
