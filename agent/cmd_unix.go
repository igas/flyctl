//go:build !windows
// +build !windows

package agent

import (
	"os/exec"
	"syscall"
)

func SetSysProcAttributes(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}
}
