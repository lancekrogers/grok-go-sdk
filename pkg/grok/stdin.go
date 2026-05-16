package grok

import (
	"bytes"
	"context"
	"io"
	"os/exec"
)

func (c *GrokClient) RunFromStdin(stdin io.Reader, prompt string, opts *RunOptions) (*GrokResult, error) {
	return c.RunFromStdinCtx(context.Background(), stdin, prompt, opts)
}

func (c *GrokClient) RunFromStdinCtx(ctx context.Context, stdin io.Reader, prompt string, opts *RunOptions) (*GrokResult, error) {
	prepared, err := c.prepareOptions(opts)
	if err != nil {
		return nil, err
	}
	if prepared.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, prepared.Timeout)
		defer cancel()
	}
	args := BuildArgs(prompt, prepared)
	cmd := execCommand(ctx, c.BinPath, args...)
	cmd.Dir = c.workDir(prepared)
	cmd.Env = c.envFor(prepared)
	cmd.Stdin = stdin

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		exitCode := 1
		if ee, ok := err.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		}
		ge := ParseError(stderr.String(), exitCode)
		ge.Original = err
		return nil, ge
	}
	return decodeOutput(prepared.Format, stdout.Bytes())
}
