package utils

import (
	"errors"
	"log"
	"os"
	"runtime"
	"syscall"
)

func Daemon(chdir bool, clstd bool) (err error) {
	var ret1, ret2 uintptr
	var er syscall.Errno

	darwin := runtime.GOOS == "darwin"
	if syscall.Getppid() == 1 {
		return
	}
	ret1, ret2, er = syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
	if er != 0 {
		err = errors.New("daemon fork fail")
		return
	}
	if ret2 < 0 {
		err = errors.New("fork return fail")
		os.Exit(-1)
	}
	if darwin && ret2 == 1 {
		ret1 = 0
	}
	if ret1 > 0 {
		os.Exit(0)
	}
	syscall.Umask(0)
	syscall.Setsid()
	if chdir {
		os.Chdir("/")
	}
	if clstd {
		fs, err := os.OpenFile(os.DevNull, os.O_RDWR, 0)
		if err == nil {
			fd := fs.Fd()
			syscall.Dup2(int(fd), int(os.Stdin.Fd()))
			syscall.Dup2(int(fd), int(os.Stdout.Fd()))
			syscall.Dup2(int(fd), int(os.Stderr.Fd()))
		}
	}
	return
}
