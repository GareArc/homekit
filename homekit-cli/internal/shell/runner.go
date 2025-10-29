package shell

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// Options configures script execution via the embedded shell interpreter.
type Options struct {
	Args    []string
	Env     map[string]string
	Dir     string
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	Timeout time.Duration
	DryRun  bool
}

// Run executes the shell script provided by reader using mvdan's interpreter.
func Run(ctx context.Context, name string, reader io.Reader, opts Options) error {
	if opts.DryRun {
		return nil
	}

	parser := syntax.NewParser()
	prog, err := parser.Parse(reader, name)
	if err != nil {
		return fmt.Errorf("parse script %s: %w", name, err)
	}

	ctx, cancel := withTimeout(ctx, opts.Timeout)
	defer cancel()

	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}
	stdin := opts.Stdin
	if stdin == nil {
		stdin = os.Stdin
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
		return fmt.Errorf("init interpreter: %w", err)
	}

	if err := runner.Run(ctx, prog); err != nil {
		return normalizeError(err)
	}
	return nil
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

func normalizeError(err error) error {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	// interpreter returns ExitStatusError with message like "exit status 1"
	if strings.HasPrefix(err.Error(), "exit status ") {
		return err
	}
	return fmt.Errorf("interpreter: %w", err)
}
