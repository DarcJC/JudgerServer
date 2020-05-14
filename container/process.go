package container

/*
#include "unistd.h"

int do_fork() {
    pid_t pid = fork();
    return pid;
}
*/
import "C"

// Fork 封装了一下原生的Fork
func Fork() int {
	return int(C.do_fork())
}
