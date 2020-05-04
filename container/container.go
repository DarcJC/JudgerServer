package container

import seccomp "github.com/seccomp/libseccomp-golang"

func init() {
}

// Container 创建一个用于运行不可信程序的容器
type Container struct {
	// Filter的生命周期从InitSeccomp开始 到LoadSeccomp结束
	Filter *seccomp.ScmpFilter

	SeccompRule []seccomp.ScmpSyscall
}

// InitSeccomp 初始化Seccomp过滤器
func (c Container) InitSeccomp() error {
	f, err := seccomp.NewFilter(seccomp.ActKill)
	if err != nil {
		return err
	}
	for _, s := range c.SeccompRule {
		f.AddRule(s, seccomp.ActAllow)
	}
	c.Filter = f
	return nil
}

// LoadSeccomp 将Seccomp过滤器载入内核
// 执行该操作会销毁Filter对象
func (c Container) LoadSeccomp() error {
	err := c.Filter.Load()
	c.Filter.Release()
	c.Filter = nil
	if err != nil {
		return err
	}
	return nil
}

// Destory 销毁并释放容器资源
func (c Container) Destory() error {
	// 确认是否已释放Filter
	if c.Filter != nil {
		c.Filter.Release()
		c.Filter = nil
	}
	return nil
}

// GetSyscallNumber 封装了一下GetSyscallFromName
// 如果有err则直接报错
func GetSyscallNumber(name string) seccomp.ScmpSyscall {
	s, err := seccomp.GetSyscallFromName(name)
	if err != nil {
		panic(err)
	}
	return s
}

// NewContainer 返回一个新的容器
// Rule 是白名单
func NewContainer(rule []seccomp.ScmpSyscall) *Container {
	if rule == nil {
		rule = []seccomp.ScmpSyscall{
			GetSyscallNumber("read"),
			GetSyscallNumber("fstat"),
			GetSyscallNumber("mmap"),
			GetSyscallNumber("mprotect"),
			GetSyscallNumber("munmap"),
			GetSyscallNumber("uname"),
			GetSyscallNumber("arch_prctl"),
			GetSyscallNumber("brk"),
			GetSyscallNumber("access"),
			GetSyscallNumber("exit_group"),
			GetSyscallNumber("close"),
			GetSyscallNumber("readlink"),
			GetSyscallNumber("sysinfo"),
			GetSyscallNumber("write"),
			GetSyscallNumber("writev"),
			GetSyscallNumber("lseek"),
			GetSyscallNumber("clock_gettime"),
		}
	}
	return &Container{
		SeccompRule: rule,
	}
}
