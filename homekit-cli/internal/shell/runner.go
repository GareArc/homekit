package shell

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

const (
	exitCodeUnknown  = -1
	exitCodeTimeout  = 124
	exitCodeCanceled = 125
)

// Options configures script execution via the embedded shell interpreter.
type Options struct {
	Args          []string
	Env           map[string]string
	Dir           string
	Stdin         io.Reader
	Stdout        io.Writer
	Stderr        io.Writer
	Timeout       time.Duration
	DryRun        bool
	CaptureOutput bool
}

// Result captures stdout/stderr and exit information from a shell script run.
type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
}

// Run executes the shell script provided by reader using mvdan's interpreter.
func Run(ctx context.Context, name string, reader io.Reader, opts Options) (Result, error) {
	res := Result{}
	if opts.DryRun {
		return res, nil
	}

	parser := syntax.NewParser()
	prog, err := parser.Parse(reader, name)
	if err != nil {
		return res, fmt.Errorf("parse script %s: %w", name, err)
	}

	ctx, cancel := withTimeout(ctx, opts.Timeout)
	defer cancel()

	stdout := opts.Stdout
	if stdout == nil {
		if opts.CaptureOutput {
			stdout = io.Discard
		} else {
			stdout = os.Stdout
		}
	}
	stderr := opts.Stderr
	if stderr == nil {
		if opts.CaptureOutput {
			stderr = io.Discard
		} else {
			stderr = os.Stderr
		}
	}
	stdin := opts.Stdin
	if stdin == nil {
		stdin = os.Stdin
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	if opts.CaptureOutput {
		stdout = io.MultiWriter(stdout, &stdoutBuf)
		stderr = io.MultiWriter(stderr, &stderrBuf)
	}

	env := buildEnv(opts.Env)

	options := []interp.RunnerOption{
		interp.Params(opts.Args...),
		interp.Env(expand.ListEnviron(env...)),
		interp.StdIO(stdin, stdout, stderr),
	}
	if opts.Dir != "" {
		options = append(options, interp.Dir(opts.Dir))
	}

	runner, err := interp.New(options...)
	if err != nil {
		return res, fmt.Errorf("init interpreter: %w", err)
	}

	start := time.Now()
	err = runner.Run(ctx, prog)
	res.Duration = time.Since(start)

	if opts.CaptureOutput {
		res.Stdout = stdoutBuf.String()
		res.Stderr = stderrBuf.String()
	}

	if err != nil {
		exitCode, normErr := normalizeError(err)
		res.ExitCode = exitCode
		return res, normErr
	}

	res.ExitCode = 0
	return res, nil
}

func buildEnv(overrides map[string]string) []string {
	base := os.Environ()
	if len(overrides) == 0 {
		return base
	}

	result := make([]string, 0, len(base)+len(overrides))
	result = append(result, base...)
	for key, value := range overrides {
		result = append(result, fmt.Sprintf("%s=%s", key, value))
	}
	return result
}

func withTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	if timeout <= 0 {
		return context.WithCancel(ctx)
	}
	return context.WithTimeout(ctx, timeout)
}

func normalizeError(err error) (int, error) {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return exitCodeTimeout, err
	case errors.Is(err, context.Canceled):
		return exitCodeCanceled, err
	}

	var status interp.ExitStatus
	if errors.As(err, &status) {
		return int(status), err
	}
	return exitCodeUnknown, fmt.Errorf("interpreter: %w", err)
}
