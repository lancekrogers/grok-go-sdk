package grok

import (
	"os/exec"
	"syscall"
	"time"
)

func stopGracefully(cmd *exec.Cmd, waitDone <-chan struct{}) {
	if cmd.Process != nil {
		_ = cmd.Process.Signal(syscall.SIGTERM)
	}
	select {
	case <-waitDone:
	case <-time.After(5 * time.Second):
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		<-waitDone
	}
}
