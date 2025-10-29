package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/homekit/homekit-cli/internal/core"
	"github.com/homekit/homekit-cli/internal/plugins"
)

// NewPluginCommand exposes plugin discovery and introspection commands.
func NewPluginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugins",
		Short: "Manage external plugin integrations",
	}

	cmd.AddCommand(newPluginListCommand())
	return cmd
}

func newPluginListCommand() *cobra.Command {
	var prefix string
	var extraPaths []string

	c := &cobra.Command{
		Use:   "list",
		Short: "List discoverable plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			rt, err := runtimeFrom(cmd)
			if err != nil {
				return err
			}

			searchPaths := append([]string{}, extraPaths...)
			searchPaths = append(searchPaths, configuredPluginPaths(rt.Config)...)
			manager := plugins.NewManager(prefixOrDefault(prefix), searchPaths)
			plugins, err := manager.Discover()
			if err != nil {
				return err
			}

			for _, plugin := range plugins {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", plugin.Name, plugin.Path)
			}
			return nil
		},
	}

	c.Flags().StringVar(&prefix, "prefix", "homekit-cli", "Plugin prefix to search for")
	c.Flags().StringSliceVar(&extraPaths, "path", nil, "Additional directories to search")
	return c
}

func prefixOrDefault(prefix string) string {
	if prefix == "" {
		return "homekit-cli"
	}
	return prefix
}

func configuredPluginPaths(cfg core.Config) []string {
	return cfg.PluginPaths
}
