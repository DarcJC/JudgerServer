#!/bin/bash
go build
sudo setcap CAP_SYS_ADMIN,CAP_SETUID,CAP_SETGID+eip ./JudgerServer

