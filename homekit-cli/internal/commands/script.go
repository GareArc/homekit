package commands

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/homekit/homekit-cli/internal/assets"
	"github.com/homekit/homekit-cli/internal/core"
	executor "github.com/homekit/homekit-cli/internal/exec"
	"github.com/homekit/homekit-cli/internal/shell"
	"github.com/homekit/homekit-cli/pkg/utils"
)

// NewScriptCommand wires the `script` command group.
func NewScriptCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "script",
		Short: "Execute local or embedded scripts",
	}

	runCmd := &cobra.Command{
		Use:   "run [path|asset]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Run a script by path or embedded name",
		RunE: func(cmd *cobra.Command, args []string) error {
			rt, err := runtimeFrom(cmd)
			if err != nil {
				return err
			}

			timeout, _ := cmd.Flags().GetDuration("timeout")
			env, _ := cmd.Flags().GetStringSlice("env")
			embeddedName, _ := cmd.Flags().GetString("embedded")
			workingDir, _ := cmd.Flags().GetString("workdir")

			spec := executor.Spec{
				Command:       args[0],
				Args:          args[1:],
				Timeout:       timeout,
				Dir:           workingDir,
				CaptureOutput: true,
				DryRun:        rt.DryRun,
				Env:           parseEnv(env),
			}

			return runScript(cmd, rt, embeddedName, spec)
		},
	}

	runCmd.Flags().String("embedded", "", "Name of embedded script to execute (overrides path)")
	runCmd.Flags().Duration("timeout", 5*time.Minute, "Timeout for the script execution")
	runCmd.Flags().StringSlice("env", nil, "Environment variables (KEY=VALUE)")
	runCmd.Flags().String("workdir", "", "Working directory for the process")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available embedded scripts",
		RunE: func(cmd *cobra.Command, args []string) error {
			rt, err := runtimeFrom(cmd)
			if err != nil {
				return err
			}
			manager := assets.NewManager(assets.Embedded(), overrideDirectory(rt.Config))
			names, err := manager.List("scripts")
			if err != nil {
				return err
			}
			for _, name := range names {
				fmt.Fprintln(cmd.OutOrStdout(), name)
			}
			return nil
		},
	}

	root.AddCommand(runCmd, listCmd)
	return root
}

func runScript(cmd *cobra.Command, rt *core.Runtime, embeddedName string, spec executor.Spec) error {
	if embeddedName == "" {
		res, err := executor.Run(cmd.Context(), spec)
		if err != nil {
			return err
		}
		if res.Stdout != "" {
			fmt.Fprint(cmd.OutOrStdout(), res.Stdout)
		}
		if res.Stderr != "" {
			fmt.Fprint(cmd.ErrOrStderr(), res.Stderr)
		}
		return nil
	}

	manager := assets.NewManager(assets.Embedded(), overrideDirectory(rt.Config))
	handle, err := manager.Open("scripts", embeddedName)
	if err != nil {
		return fmt.Errorf("open embedded script: %w", err)
	}
	defer handle.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	stdout := cmd.OutOrStdout()
	stderr := cmd.ErrOrStderr()
	if spec.CaptureOutput {
		stdout = &stdoutBuf
		stderr = &stderrBuf
	}

	if err := shell.Run(cmd.Context(), embeddedName, handle, shell.Options{
		Args:    spec.Args,
		Env:     spec.Env,
		Dir:     spec.Dir,
		Stdin:   spec.Stdin,
		Stdout:  stdout,
		Stderr:  stderr,
		Timeout: spec.Timeout,
		DryRun:  spec.DryRun,
	}); err != nil {
		return err
	}

	if spec.CaptureOutput {
		if stdoutBuf.Len() > 0 {
			fmt.Fprint(cmd.OutOrStdout(), stdoutBuf.String())
		}
		if stderrBuf.Len() > 0 {
			fmt.Fprint(cmd.ErrOrStderr(), stderrBuf.String())
		}
	}
	return nil
}

func parseEnv(values []string) map[string]string {
	env := map[string]string{}
	for _, kv := range values {
		parts := utils.SplitPair(kv, "=")
		if parts[0] == "" {
			continue
		}
		env[parts[0]] = parts[1]
	}
	return env
}

func overrideDirectory(cfg core.Config) string {
	return cfg.AssetOverrides
}

func runtimeFrom(cmd *cobra.Command) (*core.Runtime, error) {
	rt, ok := core.FromContext(cmd.Context())
	if !ok {
		return nil, errors.New("runtime unavailable")
	}
	return rt, nil
}
