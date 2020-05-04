package container

/*
#include "c_utils.h"
*/
import "C"

// Fork 封装了一下原生的Fork
func Fork() int {
	return int(C.do_fork())
}
