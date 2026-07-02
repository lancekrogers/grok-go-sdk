package grok

import (
	"errors"
	"os/exec"
	"syscall"
	"time"
)

func stopGracefully(cmd *exec.Cmd, waitDone <-chan error) error {
	select {
	case err := <-waitDone:
		return err
	default:
	}

	if cmd.Process == nil {
		return nil
	}
	signalErr := cmd.Process.Signal(syscall.SIGTERM)
	select {
	case err := <-waitDone:
		if isExpectedStopExit(err) {
			return nil
		}
		if signalErr != nil {
			return signalErr
		}
		return err
	case <-time.After(5 * time.Second):
		killErr := error(nil)
		if cmd.Process != nil {
			killErr = cmd.Process.Kill()
		}
		err := <-waitDone
		if isExpectedStopExit(err) {
			return nil
		}
		if killErr != nil {
			return killErr
		}
		return err
	}
}

func isExpectedStopExit(err error) bool {
	if err == nil {
		return true
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		return false
	}
	status, ok := exitErr.Sys().(syscall.WaitStatus)
	if !ok || !status.Signaled() {
		return false
	}
	switch status.Signal() {
	case syscall.SIGTERM, syscall.SIGKILL:
		return true
	default:
		return false
	}
}
