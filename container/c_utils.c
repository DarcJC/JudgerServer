#include "unistd.h"
#include "c_utils.h"

int do_fork() {
    pid_t pid = fork();
    return pid;
}

