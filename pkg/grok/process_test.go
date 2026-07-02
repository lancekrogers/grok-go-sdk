package grok

import (
	"os/exec"
	"testing"
	"time"
)

func TestStopGracefully_TerminatedProcessIsSuccessfulStop(t *testing.T) {
	cmd := exec.Command("sh", "-c", "while :; do sleep 1; done")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	waitDone := make(chan error, 1)
	go func() {
		waitDone <- cmd.Wait()
	}()

	if err := stopGracefully(cmd, waitDone); err != nil {
		t.Fatalf("Stop should treat SIGTERM as successful graceful stop: %v", err)
	}
}

func TestStopGracefully_ReturnsPriorProcessFailure(t *testing.T) {
	cmd := exec.Command("sh", "-c", "exit 7")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	waitDone := make(chan error, 1)
	go func() {
		waitDone <- cmd.Wait()
	}()

	deadline := time.Now().Add(2 * time.Second)
	for len(waitDone) == 0 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if len(waitDone) == 0 {
		t.Fatal("process did not exit before deadline")
	}
	if err := stopGracefully(cmd, waitDone); err == nil {
		t.Fatal("Stop should return a process error that occurred before stopping")
	}
}
