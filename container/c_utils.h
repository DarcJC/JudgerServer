#define _GNU_SOURCE

enum process_type {
    CHILD = 0,
    PARENT = 1,
    ERROR = 2,
};


int do_fork();

