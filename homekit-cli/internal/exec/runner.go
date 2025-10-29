package executor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// Spec describes a command invocation.
type Spec struct {
	Command       string
	Args          []string
	Env           map[string]string
	Dir           string
	Timeout       time.Duration
	Stdin         io.Reader
	CaptureOutput bool
	DryRun        bool
}

// Result captures execution details from a command run.
type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
}

// Run executes a command according to the specification.
func Run(ctx context.Context, spec Spec) (Result, error) {
	if spec.Command == "" {
		return Result{}, errors.New("exec: command must be specified")
	}

	res := Result{}

	if spec.DryRun {
		res.ExitCode = 0
		return res, nil
	}

	ctx, cancel := augmentContext(ctx, spec.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, spec.Command, spec.Args...)
	cmd.Dir = spec.Dir
	if spec.Env != nil {
		cmd.Env = append(os.Environ(), formatEnv(spec.Env)...)
	}
	if spec.Stdin != nil {
		cmd.Stdin = spec.Stdin
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	if spec.CaptureOutput {
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	start := time.Now()
	err := cmd.Run()
	res.Duration = time.Since(start)

	if spec.CaptureOutput {
		res.Stdout = stdoutBuf.String()
		res.Stderr = stderrBuf.String()
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			res.ExitCode = status.ExitStatus()
		} else {
			res.ExitCode = exitErr.ExitCode()
		}
		return res, fmt.Errorf("command failed: %w", err)
	}

	if err != nil {
		return res, err
	}

	res.ExitCode = 0
	return res, nil
}

func augmentContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout > 0 {
		return context.WithTimeout(ctx, timeout)
	}
	if ctx == nil {
		return context.WithCancel(context.Background())
	}
	return context.WithCancel(ctx)
}

func formatEnv(values map[string]string) []string {
	out := make([]string, 0, len(values))
	for k, v := range values {
		out = append(out, fmt.Sprintf("%s=%s", k, v))
	}
	return out
}
