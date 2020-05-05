package main

import (
	"JudgerServer/container"

	seccomp "github.com/seccomp/libseccomp-golang"
)

func init() {
}

func main() {
	// router.Router.Run(":8081")
	container.CreateRunner(&container.RunnerConfig{
		WorkDir:     "/home/darc/Code/JudgerServer/test_dir/",
		ChangeRoot:  true,
		GID:         1000,
		UID:         1000,
		RunablePath: "test",
		OutputPath:  "output",
		InputPath:   "input",
		ErrorPath:   "error",
		SeccompRule: []seccomp.ScmpSyscall{
			container.GetSyscallNumber("clone"),
			container.GetSyscallNumber("acct"),
			container.GetSyscallNumber("add_key"),
			container.GetSyscallNumber("bpf"),
			container.GetSyscallNumber("clock_adjtime"),
			container.GetSyscallNumber("clock_settime"),
			container.GetSyscallNumber("clone"),
			container.GetSyscallNumber("create_module"),
			container.GetSyscallNumber("delete_module"),
			container.GetSyscallNumber("finit_module"),
			container.GetSyscallNumber("get_kernel_syms"),
			container.GetSyscallNumber("get_mempolicy"),
			container.GetSyscallNumber("init_module"),
			container.GetSyscallNumber("ioperm"),
			container.GetSyscallNumber("iopl"),
			container.GetSyscallNumber("kcmp"),
			container.GetSyscallNumber("kexec_file_load"),
			container.GetSyscallNumber("kexec_load"),
			container.GetSyscallNumber("keyctl"),
			container.GetSyscallNumber("lookup_dcookie"),
			container.GetSyscallNumber("mbind"),
			container.GetSyscallNumber("mount"),
			container.GetSyscallNumber("move_pages"),
			container.GetSyscallNumber("name_to_handle_at"),
			container.GetSyscallNumber("nfsservctl"),
			container.GetSyscallNumber("open_by_handle_at"),
			container.GetSyscallNumber("perf_event_open"),
			container.GetSyscallNumber("personality"),
			container.GetSyscallNumber("pivot_root"),
			container.GetSyscallNumber("process_vm_readv"),
			container.GetSyscallNumber("process_vm_writev"),
			container.GetSyscallNumber("ptrace"),
			container.GetSyscallNumber("query_module"),
			container.GetSyscallNumber("quotactl"),
			container.GetSyscallNumber("reboot"),
			container.GetSyscallNumber("request_key"),
			container.GetSyscallNumber("set_mempolicy"),
			container.GetSyscallNumber("setns"),
			container.GetSyscallNumber("settimeofday"),
			container.GetSyscallNumber("stime"),
			container.GetSyscallNumber("swapon"),
			container.GetSyscallNumber("swapoff"),
			container.GetSyscallNumber("sysfs"),
			container.GetSyscallNumber("_sysctl"),
			container.GetSyscallNumber("umount"),
			container.GetSyscallNumber("umount2"),
			container.GetSyscallNumber("unshare"),
			container.GetSyscallNumber("uselib"),
			container.GetSyscallNumber("userfaultfd"),
			container.GetSyscallNumber("ustat"),
			container.GetSyscallNumber("vm86"),
			container.GetSyscallNumber("vm86old"),
		},
		SeccompType: seccomp.ActAllow,
	})
}
