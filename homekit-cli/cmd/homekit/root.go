package homekit

import (
	"context"
	"fmt"
	"sync"

	"github.com/spf13/cobra"

	"github.com/homekit/homekit-cli/internal/commands"
	"github.com/homekit/homekit-cli/internal/core"
)

var (
	version = "dev"
	commit  = "none"
	date    = ""
	source  = "unknown"
)

type rootOptions struct {
	ConfigPath string
	LogLevel   string
	LogFormat  string
	NoColor    bool
	DryRun     bool
}

var (
	opts            rootOptions
	bootstrapOnce   sync.Once
	bootstrapError  error
	runtimeInstance *core.Runtime
)

// NewRootCommand constructs the Cobra command tree for the CLI.
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "homekit",
		Short: "HomeKit CLI for orchestrating home server workflows",
		Long: `HomeKit CLI offers a plugin-friendly command surface for home server operations.
It supports executing embedded scripts, rendering templates, and delegating to external plugins.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			bootstrapOnce.Do(func() {
				ctx := cmd.Context()
				rt, err := core.Bootstrap(ctx, core.Options{
					ConfigPath: opts.ConfigPath,
					LogLevel:   opts.LogLevel,
					LogFormat:  opts.LogFormat,
					NoColor:    opts.NoColor,
					DryRun:     opts.DryRun,
				}, core.VersionInfo{
					Version: version,
					Commit:  commit,
					Date:    date,
					Source:  source,
				})
				if err != nil {
					bootstrapError = err
					return
				}
				runtimeInstance = rt
				cmd.SetContext(rt.Context)
			})
			if bootstrapError == nil && runtimeInstance != nil {
				cmd.SetContext(core.WithRuntime(cmd.Context(), runtimeInstance))
			}
			return bootstrapError
		},
	}

	cmd.PersistentFlags().StringVar(&opts.ConfigPath, "config", "", "Path to config file (default: platform config directory)")
	cmd.PersistentFlags().StringVar(&opts.LogLevel, "log-level", "info", "Log level (trace, debug, info, warn, error)")
	cmd.PersistentFlags().StringVar(&opts.LogFormat, "log-format", "console", "Log format (console|json)")
	cmd.PersistentFlags().BoolVar(&opts.NoColor, "no-color", false, "Disable ANSI colors in console output")
	cmd.PersistentFlags().BoolVar(&opts.DryRun, "dry-run", false, "Simulate actions without executing them")

	cmd.AddCommand(newVersionCommand())
	cmd.AddCommand(commands.NewScriptCommand())
	cmd.AddCommand(commands.NewAssetsCommand())
	cmd.AddCommand(commands.NewTemplateCommand())
	cmd.AddCommand(commands.NewDockerCommand())
	cmd.AddCommand(commands.NewSystemCommand())
	cmd.AddCommand(commands.NewPluginCommand())

	return cmd
}

// Execute runs the CLI root command.
func Execute() error {
	return NewRootCommand().Execute()
}

// ContextRuntime extracts the initialized runtime from command context or returns an error.
func ContextRuntime(cmd *cobra.Command) (*core.Runtime, error) {
	rt, ok := core.FromContext(cmd.Context())
	if !ok {
		return nil, fmt.Errorf("runtime not initialized")
	}
	return rt, nil
}

// SetContext injects a context with background runtime for tests.
func SetContext(cmd *cobra.Command, ctx context.Context) {
	cmd.SetContext(ctx)
}
