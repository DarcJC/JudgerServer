#!/bin/bash
go build
sudo setcap CAP_SYS_ADMIN,CAP_SETUID,CAP_SETGID+eip ./JudgerServer
sudo setcap "cap_setgid=eip cap_sys_admin=eip" ./JudgerServer

